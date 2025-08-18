package jwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	ID        string          `json:"id"`
	UserName  string          `json:"username"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	IsAdmin   string          `json:"is_admin"`
	ExpiresAr jwt.NumericDate `json:"expires_at"`
}
