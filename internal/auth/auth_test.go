package internal

import "testing"

func TestPasswordHashing(t *testing.T) {

	rightPassword := "superSecretPassword1235!"
	wrongPassword := "blinBlonBlinCarapachoWOO"

	hash, err := HashPassword(rightPassword)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned an empty string")
	}

	isValid, err := CheckPasswordHash(rightPassword, hash)
	if err != nil {
		t.Fatalf("CheckPasswordsHash encounter an unexpected error: %v", err)
	}
	if !isValid {
		t.Error("CheckPasswordHash returned false for the correct password; expected true")
	}

	isInvalid, err := CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash encountered an unexpected error on mismatch: %v", err)
	}
	if isInvalid {
		t.Errorf("CheckPasswordHash returned true for a wrong password; expected false")
	}

}