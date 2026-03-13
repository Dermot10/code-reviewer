package utils

import "testing"

func TestHashedPassword(t *testing.T) {
	password := "mysecret"

	hash, err := HashedPassword(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if hash == password {
		t.Fatalf("hash should not equal original password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "mysecret"

	hash, err := HashedPassword(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !CheckPassword(hash, password) {
		t.Errorf("expected password to match hash")
	}

	if CheckPassword(hash, "wrongpassword") {
		t.Errorf("expected password check to fail")
	}
}
