package internal

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password string, hash string) (bool, error) {
	result, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return result, err
	}
	return result, err
}
 
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer: "chirpy-access",
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			Subject: userID.String(),
		},
	)

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	
	customClaim := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&customClaim,
		func(token *jwt.Token) (interface{}, error) {return []byte(tokenSecret), nil},
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDRaw, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, fmt.Errorf("issuer is invalid")
	}
	if issuer != "chirpy-access" {
		return uuid.Nil, fmt.Errorf("issuer is invalid")
	}

	userID, err := uuid.Parse(userIDRaw)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil

}

func GetBearerToken(headers http.Header) (string, error) {
	rawHeader := headers.Get("Authorization")
	
	if rawHeader == "" {
		return "", fmt.Errorf("Unable to get header string")
	}
	
	tokenString := strings.ReplaceAll(rawHeader, "Bearer ", "")
	if tokenString == "" {
		return "", fmt.Errorf("Unable to get token string")
	}

	return tokenString, nil
}

