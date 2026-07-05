package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	App               AppConfig
	DB                DBConfig
	GoogleOauthConfig *oauth2.Config
	Cookie            CookieConfig
	JWTSecret         string
}

type CookieConfig struct {
	Name     string
	Secure   bool
	Domain   string
	HttpOnly bool
}

type AppConfig struct {
	Port string
	Env  string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "todo_app"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		GoogleOauthConfig: &oauth2.Config{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
		Cookie: CookieConfig{
			Name:     getEnv("COOKIE_NAME", "session_id"),
			Secure:   getEnv("COOKIE_SECURE", "true") == "true",
			Domain:   getEnv("COOKIE_DOMAIN", "localhost"),
			HttpOnly: getEnv("COOKIE_HTTPONLY", "true") == "true",
		},
		JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),
	}, nil
}

// key and the fallback value, if the key is not found in the env variables
func getEnv(key, defaultValue string) string {
	value := defaultValue
	if envValue, exists := os.LookupEnv(key); exists {
		value = envValue
	}
	return value
}

// data source name
func (db DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.SSLMode,
	)
}
