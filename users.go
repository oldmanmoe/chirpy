package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)



type User struct {
	ID			uuid.UUID	`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Email		string		`json:"email"`
}

type AuthRequest struct {
	Email		string	`json:"email"`
	Password	string	`json:"password"`
}

func(cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request,) {
	var registerReq AuthRequest

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&registerReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	password, err := internal.HashPassword(registerReq.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong hashing password")
	}

	userInfo, err := cfg.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			Email: registerReq.Email,
			HashedPassword: password,
		})
	if err != nil {
		log.Fatal(err)
		return 
	}

	newUser := User{
		ID: userInfo.ID,
		CreatedAt: userInfo.CreatedAt,
		UpdatedAt: userInfo.UpdatedAt,
		Email: userInfo.Email,
	}

	printUserInfo(newUser)
	
	respondWithJSON(w, http.StatusCreated, newUser)
}

func(cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, req *http.Request) {
	var loginReq AuthRequest

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&loginReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	existingUser, err := cfg.db.GetUserByEmail(context.Background(), loginReq.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	isValid, err := internal.CheckPasswordHash(loginReq.Password, existingUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	
	if !isValid {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	resultUser := User{
		ID: existingUser.ID,
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: existingUser.UpdatedAt,
		Email: existingUser.Email,
	}

	respondWithJSON(w, http.StatusOK, resultUser)

}

func resetUsers(apiConfig *apiConfig) {
	if err := apiConfig.db.ResetUsers(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func printUserInfo(user User) string {
	return fmt.Sprintf(`
ID:.......... %v
CreatedAt.... %v
UpdatedAt.... %v
Email........ %v`, user.ID, user.CreatedAt, user.UpdatedAt, user.Email)
}