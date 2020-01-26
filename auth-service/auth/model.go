package auth

import "time"

type (
	TokenRequestBody struct {
		UserID string `'json:userId`
	}
	TokenInfo struct {
		AccessToken  AccessToken `json:"access_token"`
		RefreshToken AccessToken `json:"refresh_token"`
		ExprireAt    int64       `json:"exp"`
	}

	AccessToken struct {
		Audience  string   `json:"aud"`
		Issuer    string   `json:"iss"`
		JTI       string   `json:"jti"`
		Sub       string   `json:"sub"`
		Roles     []string `json:"roles"`
		ExprireAt int64    `json:"exp"`
	}

	CreateUserRequest struct {
		UserID   string `json:"uid"`
		Name     string `json:"name"`
		Email    string
		Root_dir string `json:"root_dir"`
	}
	CreateUserResponse struct {
		Ok              bool `json:"ok"`
		IsAuthenticated bool `json:"is_authenticated"`
		Token           string
		ExprireAt       time.Time
	}
	GetUserRequest struct {
		Id string `json:"id"`
	}
	GetUserResponse struct {
		Ok   bool `json:"ok"`
		Data User `json:"user"`
	}
)

type User struct {
	UId      string
	Name     string
	Root_dir string
	//Email string
}
