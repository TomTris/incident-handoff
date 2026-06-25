package main

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func mustHash(t *testing.T, pw string) string {
	t.Helper()
	h, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword(%q) returned error: %v", pw, err)
	}
	if h == "" {
		t.Fatalf("HashPassword(%q) returned empty hash", pw)
	}
	return h
}

func TestPasswordRoundtrip(t *testing.T) {
	t.Run("correct password verifies", func(t *testing.T) {
		h := mustHash(t, "correct-password")
		if err := VerifyPassword(h, "correct-password"); err != nil {
			t.Fatalf("expected match, got %v", err)
		}
	})

	t.Run("wrong password is rejected", func(t *testing.T) {
		h := mustHash(t, "correct-password")
		err := VerifyPassword(h, "wrong-password")
		if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			t.Fatalf("expected ErrMismatchedHashAndPassword, got %v", err)
		}
	})

	t.Run("salt makes hashes non-deterministic", func(t *testing.T) {
		h1 := mustHash(t, "correct-password")
		h2 := mustHash(t, "correct-password")
		if h1 == h2 {
			t.Fatal("expected differing hashes from per-call salt")
		}
		// both must still verify despite differing
		if err := VerifyPassword(h2, "correct-password"); err != nil {
			t.Fatalf("second hash failed to verify: %v", err)
		}
	})

	// The reason this construction exists: bcrypt alone caps at 72 bytes.
	t.Run("password over 72 bytes verifies", func(t *testing.T) {
		long := strings.Repeat("a", 100) // 100 bytes, exceeds bcrypt's 72-byte limit
		h := mustHash(t, long)
		if err := VerifyPassword(h, long); err != nil {
			t.Fatalf("100-byte password failed to verify: %v", err)
		}
	})

	// Distinct passwords sharing a 72-byte prefix must NOT collide.
	// Without the SHA-256 pre-hash, bcrypt truncates and these would match.
	t.Run("long passwords sharing 72-byte prefix do not collide", func(t *testing.T) {
		prefix := strings.Repeat("a", 72)
		pwX := prefix + "X"
		pwY := prefix + "Y"
		h := mustHash(t, pwX)
		if err := VerifyPassword(h, pwY); err == nil {
			t.Fatal("passwords differing only past byte 72 collided — truncation bug")
		}
	})

	// Multi-byte input: the base64 step exists to guarantee no NUL reaches bcrypt.
	t.Run("multibyte password verifies", func(t *testing.T) {
		pw := strings.Repeat("ế", 40) // 120 bytes in UTF-8, multi-byte runes
		h := mustHash(t, pw)
		if err := VerifyPassword(h, pw); err != nil {
			t.Fatalf("multibyte password failed to verify: %v", err)
		}
		if err := VerifyPassword(h, strings.Repeat("ế", 39)); err == nil {
			t.Fatal("shorter multibyte password incorrectly matched")
		}
	})

	t.Run("empty password roundtrips", func(t *testing.T) {
		h := mustHash(t, "")
		if err := VerifyPassword(h, ""); err != nil {
			t.Fatalf("empty password failed to verify: %v", err)
		}
		if err := VerifyPassword(h, "x"); err == nil {
			t.Fatal("non-empty matched empty hash")
		}
	})

	t.Run("malformed stored hash is not a mismatch error", func(t *testing.T) {
		err := VerifyPassword("not-a-bcrypt-hash", "anything")
		if err == nil {
			t.Fatal("expected error for malformed hash")
		}
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			t.Fatal("malformed hash must not report as wrong-password; it is an integrity error")
		}
	})
}

func TestIssueToken(t *testing.T) {
	user := User{
		ID:       "123",
		Username: "anh",
		Role:     "engineer",
	}
	secret := "random-jwt-secret"
	now := time.Now()
	ttl := time.Duration(15 * time.Minute)

	t.Run("normal behavior", func(t *testing.T) {
		tokenSigned, err := IssueToken(user, []byte(secret), ttl, now)
		if err != nil {
			t.Fatalf("expect no error, get error %v", err.Error())
		}

		var claims CustomClaims
		token, err := jwt.ParseWithClaims(tokenSigned, &claims, func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || token.Valid == false {
			msg := "jwt invalid"
			if err != nil {
				msg = "error: " + err.Error()
			}
			t.Fatalf("expect normal jwt decode, but %v", msg)
		}

		if claims.Subject != user.ID {
			t.Errorf("not the same ID")
		}
		if claims.Username != user.Username {
			t.Errorf("not the same Username")
		}
		if claims.Role != user.Role {
			t.Errorf("not the same Role")
		}
		if claims.IssuedAt.Unix() != now.Unix() {
			t.Errorf("IssuedAt %v", claims.IssuedAt)
		}
		if claims.ExpiresAt.Unix() != (now.Add(15 * time.Minute)).Unix() {
			t.Errorf("ExpiresAt %v", claims.ExpiresAt)
		}
		if claims.ID == "" {
			t.Errorf("JWT ID empty")
		}
	})
	t.Run("not correct token", func(t *testing.T) {
		tokenSigned, err := IssueToken(user, []byte(secret), ttl, now)
		if err != nil {
			t.Fatalf("expect no error, get error %v", err.Error())
		}

		var claims CustomClaims
		_, err = jwt.ParseWithClaims(tokenSigned, &claims, func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}
			return []byte("wrong-jwt-secret"), nil
		})

		if err == nil {
			t.Errorf("wrong secret should fail")
		}
	})
	t.Run("test using Algorithm", func(t *testing.T) {
		tokenSigned, _ := IssueToken(user, []byte(secret), ttl, now)
		parsed, _ := jwt.Parse(tokenSigned, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if parsed.Method.Alg() != "HS256" {
			t.Errorf("alg expect HS256, get %v", parsed.Method.Alg())
		}
	})
}
