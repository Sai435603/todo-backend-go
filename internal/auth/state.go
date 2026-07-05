package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/Sai435603/todo-backend-go/internal/config"
)

func GenerateState(w http.ResponseWriter, cfg config.CookieConfig) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	state := base64.URLEncoding.EncodeToString(b)

	cookie := &http.Cookie{
		Name:     cfg.Name,
		Value:    state,
		MaxAge:   int((10 * time.Minute).Seconds()),
		Secure:   cfg.Secure,
		Domain:   cfg.Domain, 
		HttpOnly: cfg.HttpOnly,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
	return state, nil
}

func GetStateCookie(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func ClearStateCookie(w http.ResponseWriter, cfg config.CookieConfig) {
	cookie := &http.Cookie{
		Name:     cfg.Name,
		Value:    "",
		MaxAge:   -1, 
		Secure:   cfg.Secure,
		Domain:   cfg.Domain,
		HttpOnly: cfg.HttpOnly,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}
