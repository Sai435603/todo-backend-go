package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

// GoogleUser represents the profile data returned by Google's userinfo endpoint.
type GoogleUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"picture"`
}

// FetchGoogleUser calls Google's userinfo API with the given access token
// and returns the authenticated user's profile.
func FetchGoogleUser(accessToken string) (*GoogleUser, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest(http.MethodGet, googleUserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("google userinfo request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo returned status %d", resp.StatusCode)
	}

	var user GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decode google userinfo: %w", err)
	}

	if user.ID == "" || user.Email == "" {
		return nil, fmt.Errorf("incomplete google user profile")
	}

	return &user, nil
}
