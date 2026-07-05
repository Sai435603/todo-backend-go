package auth

import (
	"net/http"

	"github.com/Sai435603/todo-backend-go/internal/config"
	"golang.org/x/oauth2"
)

// OAuthService encapsulates OAuth-related operations so that handlers
// never need to know about config details.
type OAuthService struct {
	oauthCfg  *oauth2.Config
	cookieCfg config.CookieConfig
}

// NewOAuthService creates an OAuthService from the relevant config slices.
// This is the only place where config is unpacked for OAuth.
func NewOAuthService(oauthCfg *oauth2.Config, cookieCfg config.CookieConfig) *OAuthService {
	return &OAuthService{
		oauthCfg:  oauthCfg,
		cookieCfg: cookieCfg,
	}
}

// AuthCodeURL returns the Google consent screen URL with the given state.
func (s *OAuthService) AuthCodeURL(state string) string {
	return s.oauthCfg.AuthCodeURL(state)
}

// GenerateState creates a cryptographically random state token and sets it
// as an HTTP cookie for CSRF protection.
func (s *OAuthService) GenerateState(w http.ResponseWriter) (string, error) {
	return GenerateState(w, s.cookieCfg)
}

// GetStateCookie reads the OAuth state value from the request cookie.
func (s *OAuthService) GetStateCookie(r *http.Request) (string, error) {
	return GetStateCookie(r, s.cookieCfg.Name)
}

// ClearStateCookie removes the OAuth state cookie from the client.
func (s *OAuthService) ClearStateCookie(w http.ResponseWriter) {
	ClearStateCookie(w, s.cookieCfg)
}

// Exchange trades an authorization code for an OAuth2 token.
func (s *OAuthService) Exchange(r *http.Request, code string) (*oauth2.Token, error) {
	return s.oauthCfg.Exchange(r.Context(), code)
}
