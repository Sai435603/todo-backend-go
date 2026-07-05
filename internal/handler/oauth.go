package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sai435603/todo-backend-go/internal/auth"
	"github.com/Sai435603/todo-backend-go/internal/config"
	"github.com/Sai435603/todo-backend-go/internal/validator"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	AuthConfig *oauth2.Config
	Cookie     config.CookieConfig
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
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Google returned an invalid status code", http.StatusInternalServerError)
		return
	}
	var userInfo map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to parse user profile", http.StatusInternalServerError)
		return
	}
	// fmt.Printf("\n\n\n\n")
	// fmt.Println(userInfo)
	// fmt.Printf("\n\n\n\n")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
