package auth

import (
	
	"testing"
	"time"

	"github.com/google/uuid"
)


/*
// go
kf := func(t *jwt.Token) (interface{}, error) {
    if t.Method != jwt.SigningMethodHS256 {
        return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
    }
    return []byte(tokenSecret), nil
}

claims := &jwt.RegisteredClaims{}
token, err := jwt.ParseWithClaims(tokenString, claims, kf)
*/

// go
func TestMakeAndValidateJWT_Success(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	tok, err := MakeJWT(userID, secret, time.Minute)
	if err != nil || tok == "" {
		t.Fatalf("expected token, got err=%v tok=%q", err, tok)
	}

	gotID, err := ValidateJWT(tok, secret)
	if err != nil {
		t.Fatalf("validate error: %v", err)
	}
	if gotID != userID {
		t.Fatalf("want %v, got %v", userID, gotID)
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	tok, _ := MakeJWT(userID, "right", time.Minute)

	_, err := ValidateJWT(tok, "wrong")
	if err == nil {
		t.Fatal("expected error with wrong secret")
	}
}

func TestValidateJWT_Expired(t *testing.T) {
	userID := uuid.New()
	tok, _ := MakeJWT(userID, "secret", time.Millisecond*50)
	time.Sleep(time.Millisecond * 60)

	_, err := ValidateJWT(tok, "secret")
	if err == nil {
		t.Fatal("expected expiration error")
	}
}

func TestValidateJWT_Malformed(t *testing.T) {
	_, err := ValidateJWT("not-a-jwt", "secret")
	if err == nil {
		t.Fatal("expected error for malformed token")
	}
}
