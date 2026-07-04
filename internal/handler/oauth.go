package handler

import (
	"net/http"

	"github.com/Sai435603/todo-backend-go/internal/auth"
	"github.com/Sai435603/todo-backend-go/internal/config"
	"github.com/Sai435603/todo-backend-go/internal/validator"
)

type AuthHandler struct {
	Config *config.Config
}

func (h *AuthHandler) HandleOAuthLogin(w http.ResponseWriter, r *http.Request) {
	state, err := auth.GenerateState(w, h.Config.Cookie)
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	url := h.Config.GoogleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	cookieState, err := auth.GetStateCookie(r, h.Config.Cookie.Name)
	if err != nil {
		http.Error(w, "State cookie missing", http.StatusBadRequest)
		return
	}

	urlState := r.URL.Query().Get("state")

	if !validator.OAuthState(cookieState, urlState) {
		http.Error(w, "Invalid state parameter", http.StatusForbidden)
		return
	}
	auth.ClearStateCookie(w, h.Config.Cookie)
}
