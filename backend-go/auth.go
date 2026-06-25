package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`    // never returned in JSON
	Role     string `json:"role"` // "engineer" or "admin"
}

type UserContext struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CustomClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

func HashPassword(password string) (string, error) {
	sum := sha256.Sum256([]byte(password))
	encoded := base64.StdEncoding.EncodeToString(sum[:])
	hash, err := bcrypt.GenerateFromPassword([]byte(encoded), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(storedHash, plain string) error {
	sum := sha256.Sum256([]byte(plain))
	encoded := base64.StdEncoding.EncodeToString(sum[:])
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(encoded))
}

func IssueToken(user User, secret []byte, ttl time.Duration, now time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
		Username: user.Username,
		Role:     user.Role,
	})
	return token.SignedString(secret)
}

type UserRegistration struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	usernameMinLen = 3
	usernameMaxLen = 15
	passwordMinLen = 8
	passwordMaxLen = 72 // bcrypt byte ceiling
)

func (u *UserRegistration) Validate() error {
	u.Username = strings.TrimSpace(u.Username) // normalize what you store

	switch n := utf8.RuneCountInString(u.Username); {
	case n == 0:
		return errors.New("username cannot be empty")
	case n < usernameMinLen:
		return fmt.Errorf("username must be at least %d characters", usernameMinLen)
	case n > usernameMaxLen:
		return fmt.Errorf("username must be at most %d characters", usernameMaxLen)
	}

	switch {
	case utf8.RuneCountInString(u.Password) < passwordMinLen:
		return fmt.Errorf("password must be at least %d characters", passwordMinLen)
	case len(u.Password) > passwordMaxLen:
		return fmt.Errorf("password must be at most %d bytes", passwordMaxLen)
	}
	return nil
}
