package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

func CreateToken(userID string) (*Auth, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    jwt.NewNumericDate(time.Now().Add(time.Minute * 30)).Unix(),
		"iat":    jwt.NewNumericDate(time.Now()).Unix(),
	}
	reflashclaims := jwt.MapClaims{
		"userID": userID,
		"exp":    jwt.NewNumericDate(time.Now().AddDate(0, 0, 7)).Unix(),
		"iat":    jwt.NewNumericDate(time.Now()).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	reflashToken := jwt.NewWithClaims(jwt.SigningMethodHS256, reflashclaims)
	signedToken, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return nil, err
	}
	signedReflashToken, err := reflashToken.SignedString([]byte("your-secret-key"))
	if err != nil {
		return nil, err
	}
	auth := &Auth{
		Token:        signedToken,
		RefreshToken: signedReflashToken,
	}
	return auth, nil
}

func ReflashToken(token string) (string, error) {
	parsedToken, err := ParseToken(token)
	if err != nil {
		return "", err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userID := claims["userID"].(string)
		newToken, err := CreateToken(userID)
		if err != nil {
			return "", err
		}
		return newToken.Token, nil
	}
	return "", errors.New("invalid token")
}

func ParseToken(jwtToken string) (*jwt.Token, error) {
	if jwtToken == "" {
		return nil, errors.New("token is empty")
	}

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte("your-secret-key"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil

}
