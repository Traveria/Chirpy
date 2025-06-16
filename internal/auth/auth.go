package auth

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := jwt.NewNumericDate(time.Now())
	futureTime := now.Add(time.Hour + expiresIn)
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  now,
		ExpiresAt: jwt.NewNumericDate(futureTime),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	finalToken, err := token.SignedString(tokenSecret)
	if err != nil {
		log.Printf("unable to make token %s", err)
		return "", err
	}

	return finalToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	tkStr, err := jwt.ParseWithClaims(tokenString, jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		tokenByte := []byte(tokenSecret)
		if t.Method != jwt.SigningMethodHS256 {
			log.Printf("incorrect signing method %s", t.Method.Alg())
			return nil, errors.New("incorrect singing method")
		}
		return tokenByte, nil
	})
	if err != nil {
		log.Printf("unable to parse token %s", err)
		return uuid.Nil, err
	}

	claims, ok := tkStr.Claims.(*jwt.RegisteredClaims)
	if !ok {
		log.Printf("cannot access claims.subject")
		return uuid.Nil, errors.New("cannot access claims.Subject")
	}

	userid, err := uuid.Parse(claims.Subject)
	if err != nil {
		log.Printf("uuid not able to convert to uuid")
		return uuid.Nil, errors.New("UUID invalid")
	}
	return userid, nil

}
