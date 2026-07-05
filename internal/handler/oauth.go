package handler

import (
	"net/http"

	"github.com/Sai435603/todo-backend-go/internal/auth"
	"github.com/Sai435603/todo-backend-go/internal/validator"
)

// AuthHandler handles OAuth-related HTTP requests.
// It depends on OAuthService, not on config — handlers should never
// know about application configuration details.
type AuthHandler struct {
	oauth *auth.OAuthService
}

// NewAuthHandler creates an AuthHandler with the given OAuthService.
func NewAuthHandler(oauth *auth.OAuthService) *AuthHandler {
	return &AuthHandler{oauth: oauth}
}

// HandleOAuthLogin initiates the Google OAuth flow by redirecting
// the user to Google's consent screen with a CSRF state token.
func (h *AuthHandler) HandleOAuthLogin(w http.ResponseWriter, r *http.Request) {
	state, err := h.oauth.GenerateState(w)
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	url := h.oauth.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// OAuthCallbackHandler handles the OAuth callback from Google.
// It validates the state parameter to prevent CSRF attacks, then
// clears the state cookie.
func (h *AuthHandler) OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	cookieState, err := h.oauth.GetStateCookie(r)
	if err != nil {
		http.Error(w, "State cookie missing", http.StatusBadRequest)
		return
	}

	urlState := r.URL.Query().Get("state")

	if !validator.OAuthState(cookieState, urlState) {
		http.Error(w, "Invalid state parameter", http.StatusForbidden)
		return
	}
	h.oauth.ClearStateCookie(w)
}
