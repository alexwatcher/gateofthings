package jwt

import (
	"time"

	"github.com/alexwatcher/gateofthings/auth/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, tokenTTL time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	clainms := token.Claims.(jwt.MapClaims)
	clainms["uid"] = user.ID
	clainms["email"] = user.Email
	clainms["exp"] = time.Now().Add(tokenTTL).Unix()

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
