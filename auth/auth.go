package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("rahasia")

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

// GENERATE JWT Token
func GenerateJWT(email, username string) (tokenString string, err error) {
	expTime := time.Now().Add(1 * time.Hour)

	claims := &JWTClaim{
		Email:    email,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)

	return
}

// Validate Token
func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return
	}

	// fmt.Println("token :", token)

	claims, ok := token.Claims.(*JWTClaim)

	// fmt.Println("claims :", claims)

	// if claims failed
	if !ok {
		err = errors.New("Could not parse claims from token")
		return
	}

	// if token expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("Token expired")
		return
	}

	return
}
