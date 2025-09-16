package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	fmt.Printf("==> Password: [%s]\n", password)
	cost := 12
	// cost, err := bcrypt.Cost([]byte(password)) // bcrypt.GenerateFromPassword(int(cost))
	// if err != nil {
	// 	return "", err
	// }

	fmt.Printf("==> Password: [%s]\tcost: [%d]\n", password, cost)

	// Hash the password using the salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	fmt.Printf("==> hashedPassword: [%s]\tcost: [%d]\n", hashedPassword, cost)
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	// Hash the provided password using the same salt as the stored hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// go
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
    claims := &jwt.RegisteredClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
        if t.Method != jwt.SigningMethodHS256 {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return []byte(tokenSecret), nil
    })
    if err != nil {
        return uuid.Nil, err
    }
    if !token.Valid {
        return uuid.Nil, errors.New("invalid token")
    }
    id, err := uuid.Parse(claims.Subject)
    if err != nil {
        return uuid.Nil, fmt.Errorf("invalid user ID in claims: %w", err)
    }
    return id, nil
}


// GetBearerToken extracts the bearer token from the Authorization header.
// It returns the token string or an error if the header is missing or malformed.
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return parts[1], nil
}