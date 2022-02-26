package security

import (
	"TeachAssistApi/app"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

type jwtClaims struct {
	Username      string `json:"usr"`
	StudentID     string `json:"st"`
	Notifications bool   `json:"nt"`
	jwt.RegisteredClaims
}

type ParsedToken struct {
	Valid         bool
	Username      string
	StudentID     string
	Notifications bool
}

func CreateJWT(username, id string, notifications bool) (string, error) {
	key := []byte(os.Getenv("JWT_KEY"))

	claims := jwtClaims{
		Username:      username,
		StudentID:     id,
		Notifications: notifications,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ta_api",
			Subject:   fmt.Sprintf("%s(%s)", username, id),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(key)
	if err != nil {
		return "", app.CreateError(app.CryptographyError)
	}

	return ss, nil
}

func VerifyJWT(tokenString string) ParsedToken {
	key := []byte(os.Getenv("JWT_KEY"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return ParsedToken{Valid: false}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		err := claims.Valid()
		if err != nil {
			return ParsedToken{Valid: false}
		}
		if !(claims["iss"] == "ta_api" && fmt.Sprintf("%s(%s)", claims["usr"], claims["st"]) == claims["sub"]) {
			return ParsedToken{Valid: false}
		}
		username, uOk := claims["usr"].(string)
		studentId, sOk := claims["st"].(string)
		notifications, nOk := claims["nt"].(bool)

		if !(uOk && sOk && nOk) {
			return ParsedToken{Valid: false}
		}

		return ParsedToken{
			Valid:         true,
			Username:      username,
			StudentID:     studentId,
			Notifications: notifications,
		}
	} else {
		return ParsedToken{Valid: false}
	}
}
