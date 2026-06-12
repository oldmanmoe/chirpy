package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

// run test with: go test -coverprofile=c.out
func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
    tests := map[string]struct {
        headers   http.Header
        wantToken string
        wantErr   bool  // true if we expect a non-nil error
    }{
        "Valid bearer token": {
            headers:   http.Header{"Authorization": []string{"Bearer abc123"}},
            wantToken: "abc123",
            wantErr:   false,
        },
        "Missing Authorization header": {
            headers:   http.Header{},
            wantToken: "",
            wantErr:   true,
        },
        "Empty Authorization header": {
            headers:   http.Header{"Authorization": []string{""}},
            wantToken: "",
            wantErr:   true,
        },
        "Bearer with no token": {
            headers:   http.Header{"Authorization": []string{"Bearer "}},
            wantToken: "",
            wantErr:   true,
        },
		"Regular Bearer 2": {
			headers: 	http.Header{"Authorization": []string{"Bearer abc.123.xyz"}},
			wantToken: 	"abc.123.xyz",
			wantErr:	false,
		},
    }

    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            gotToken, err := GetBearerToken(tc.headers)

            if tc.wantErr && err == nil {
                t.Fatalf("%s: expected an error but got none", name)
            }
            if !tc.wantErr && err != nil {
                t.Fatalf("%s: expected no error but got: %v", name, err)
            }
            if gotToken != tc.wantToken {
                t.Fatalf("%s: expected token %q, got %q", name, tc.wantToken, gotToken)
            }
        })
    }
}



