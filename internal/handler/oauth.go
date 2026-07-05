package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sai435603/todo-backend-go/internal/auth"
	"github.com/Sai435603/todo-backend-go/internal/config"
	"github.com/Sai435603/todo-backend-go/internal/database/sqlc"
	"github.com/Sai435603/todo-backend-go/internal/validator"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	AuthConfig *oauth2.Config
	Cookie     config.CookieConfig
	JWTSecret  string
	Queries    *sqlc.Queries
}

func (h *AuthHandler) HandleOAuthLogin(w http.ResponseWriter, r *http.Request) {
	state, err := auth.GenerateState(w, h.Cookie)
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	url := h.AuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	cookieState, err := auth.GetStateCookie(r, h.Cookie.Name)
	if err != nil {
		http.Error(w, "State cookie missing", http.StatusBadRequest)
		return
	}

	urlState := r.URL.Query().Get("state")
	urlCode := r.URL.Query().Get("code")
	if !validator.OAuthState(cookieState, urlState) {
		http.Error(w, "Invalid state parameter", http.StatusForbidden)
		return
	}
	auth.ClearStateCookie(w, h.Cookie)
	token, err := h.AuthConfig.Exchange(context.Background(), urlCode)
	if err != nil {
		http.Error(w, "Failed to exchange authorization code", http.StatusInternalServerError)
		return
	}

	client := h.AuthConfig.Client(r.Context(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to parse user profile", http.StatusInternalServerError)
		return
	}

	// Upsert user into database
	dbUser, err := h.Queries.UpsertUser(r.Context(), sqlc.UpsertUserParams{
		GoogleID: userInfo.ID,
		Email:    userInfo.Email,
		Name:     userInfo.Name,
	})
	if err != nil {
		fmt.Println("Failed to upsert user:", err)
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	// Generate local JWT
	jwtToken, err := auth.GenerateJWT(dbUser.ID, dbUser.Email, dbUser.Name, h.JWTSecret)
	if err != nil {
		http.Error(w, "Failed to generate session token", http.StatusInternalServerError)
		return
	}

	// Set JWT as HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		Path:     "/",
		MaxAge:   int((24 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   h.Cookie.Secure,
		Domain:   h.Cookie.Domain,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// HandleAuthCheck returns the current user info if authenticated, or 401 if not.
func (h *AuthHandler) HandleAuthCheck(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, `{"authenticated":false}`, http.StatusUnauthorized)
		return
	}

	claims, err := auth.ParseJWT(cookie.Value, h.JWTSecret)
	if err != nil {
		http.Error(w, `{"authenticated":false}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"authenticated": true,
		"user": map[string]any{
			"id":    claims.UserID,
			"email": claims.Email,
			"name":  claims.Name,
		},
	})
}

// HandleLogout clears the auth_token cookie.
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.Cookie.Secure,
		Domain:   h.Cookie.Domain,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "logged out"})
}
