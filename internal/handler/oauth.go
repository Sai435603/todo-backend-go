package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Sai435603/todo-backend-go/internal/auth"
	"github.com/Sai435603/todo-backend-go/internal/response"
	"github.com/Sai435603/todo-backend-go/internal/service"
	"github.com/Sai435603/todo-backend-go/internal/validator"
)

// AuthHandler handles OAuth-related HTTP requests.
// It depends on OAuthService, not on config — handlers should never
// know about application configuration details.
type AuthHandler struct {
	oauth   *auth.OAuthService
	userSvc *service.UserService
	jwt     *auth.JWTService
}

// NewAuthHandler creates an AuthHandler with the given dependencies.
func NewAuthHandler(oauth *auth.OAuthService, userSvc *service.UserService, jwt *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		oauth:   oauth,
		userSvc: userSvc,
		jwt:     jwt,
	}
}

// AuthResponse is the JSON payload returned after successful OAuth login.
type AuthResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// UserInfo is the public-facing subset of user data.
type UserInfo struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// HandleOAuthLogin initiates the Google OAuth flow by redirecting
// the user to Google's consent screen with a CSRF state token.
func (h *AuthHandler) HandleOAuthLogin(w http.ResponseWriter, r *http.Request) {
	state, err := h.oauth.GenerateState(w)
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	authURL := h.oauth.AuthCodeURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// OAuthCallbackHandler handles the OAuth callback from Google.
// After validating and exchanging the code, it redirects the browser
// back to the frontend with the JWT and user info encoded in the URL
// fragment (hash). This avoids COOP issues with popups.
func (h *AuthHandler) OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Validate CSRF state
	cookieState, err := h.oauth.GetStateCookie(r)
	if err != nil {
		http.Redirect(w, r, "/?auth_error=state_missing", http.StatusTemporaryRedirect)
		return
	}

	urlState := r.URL.Query().Get("state")
	if !validator.OAuthState(cookieState, urlState) {
		http.Redirect(w, r, "/?auth_error=invalid_state", http.StatusTemporaryRedirect)
		return
	}
	h.oauth.ClearStateCookie(w)

	// 2. Exchange authorization code for token
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, "/?auth_error=missing_code", http.StatusTemporaryRedirect)
		return
	}

	oauthToken, err := h.oauth.Exchange(r, code)
	if err != nil {
		http.Redirect(w, r, "/?auth_error=exchange_failed", http.StatusTemporaryRedirect)
		return
	}

	// 3. Fetch Google user profile
	googleUser, err := auth.FetchGoogleUser(oauthToken.AccessToken)
	if err != nil {
		http.Redirect(w, r, "/?auth_error=profile_failed", http.StatusTemporaryRedirect)
		return
	}

	// 4. Upsert user in database
	user, err := h.userSvc.FindOrCreateUser(r.Context(), googleUser)
	if err != nil {
		http.Redirect(w, r, "/?auth_error=save_failed", http.StatusTemporaryRedirect)
		return
	}

	// 5. Generate JWT
	token, err := h.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		http.Redirect(w, r, "/?auth_error=token_failed", http.StatusTemporaryRedirect)
		return
	}

	// 6. Encode user info as JSON and redirect to frontend with data in fragment.
	//    The fragment (#) is never sent to the server, keeping the token client-side only.
	userInfo := UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		AvatarURL: user.AvatarUrl.String,
	}
	userJSON, _ := json.Marshal(userInfo)

	redirectURL := "/#token=" + url.QueryEscape(token) + "&user=" + url.QueryEscape(string(userJSON))
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// HandleGetMe returns the currently authenticated user's profile.
func (h *AuthHandler) HandleGetMe(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	user, err := h.userSvc.GetUser(r.Context(), userID)
	if err != nil {
		_ = response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	_ = response.JSON(w, http.StatusOK, UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		AvatarURL: user.AvatarUrl.String,
	})
}
