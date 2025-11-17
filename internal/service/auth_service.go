// auth_service.go
package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type AuthService struct {
	authURL string
}

type AuthUser struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	Login       string   `json:"login"`
	Enabled     bool     `json:"enabled"`
}

// Crea un nuevo AuthService
func NewAuthService() *AuthService {
	return &AuthService{
		authURL: os.Getenv("AUTH_SERVICE_URL"),
	}
}

func (a *AuthService) IsAdmin(user *AuthUser) bool {
	for _, perm := range user.Permissions {
		if perm == "admin" {
			return true
		}
	}
	return false
}

// Valida el token JWT llamando al microservicio de autenticación
func (a *AuthService) ValidateToken(token string) (*AuthUser, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/current", a.authURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid token")
	}

	var user AuthUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, errors.New("user disabled")
	}

	fmt.Println("AUTH_SERVICE_URL =", a.authURL)
	fmt.Println("Calling:", fmt.Sprintf("%s/users/current", a.authURL))

	return &user, nil
}
