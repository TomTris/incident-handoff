package main

import (
	"context"
	"errors"
	"testing"
)

func TestGetByUsername(t *testing.T) {
	pwd1, err1 := HashPassword("anh123")
	pwd2, err2 := HashPassword("bernd123")
	pwd3, err3 := HashPassword("admin123")
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("HashPassword has problem")
	}

	var seedUsers = []User{
		{ID: "u1", Username: "anh", Password: pwd1, Role: "engineer"},
		{ID: "u2", Username: "bernd", Password: pwd2, Role: "engineer"},
		{ID: "u3", Username: "admin", Password: pwd3, Role: "admin"},
	}

	users := NewMemoryUserStoreWithSeed(seedUsers)
	for _, each := range users.users {
		_, err := users.GetByUsername(context.Background(), each.Username)
		if err != nil {
			t.Fatalf("expect no error")
		}
	}

	_, err := users.GetByUsername(context.Background(), "not-exist-user")
	if errors.Is(err, ErrUserNotFound) == false {
		t.Fatalf("expect ErrUserNotFound")
	}
}

func TestCreateUser(t *testing.T) {

	pwd1, err1 := HashPassword("anh123")
	pwd2, err2 := HashPassword("bernd123")
	pwd3, err3 := HashPassword("admin123")
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("HashPassword has problem")
	}

	var seedUsers = []User{
		{ID: "u1", Username: "anh", Password: pwd1, Role: "engineer"},
		{ID: "u2", Username: "bernd", Password: pwd2, Role: "engineer"},
		{ID: "u3", Username: "admin", Password: pwd3, Role: "admin"},
	}

	users := NewMemoryUserStoreWithSeed([]User{})

	t.Run("normal creation with sequential IDs", func(t *testing.T) {
		u0, err0 := users.Create(t.Context(), seedUsers[0])
		u1, err1 := users.Create(t.Context(), seedUsers[1])
		if err0 != nil || err1 != nil {
			t.Fatalf("expect no error")
		}
		if u0.ID != UserPrefix+"1" || u1.ID != UserPrefix+"2" {
			t.Fatalf("userID not as expected")
		}
	})

	t.Run("create an user with an existing username", func(t *testing.T) {
		_, err := users.Create(t.Context(), seedUsers[1])
		if errors.Is(err, ErrUserAlreadyExist) == false {
			t.Fatalf("expect error `%v`, get `%v`", ErrUserAlreadyExist, err)
		}
	})

}
