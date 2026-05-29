package main

import (
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

type emailRequest struct {
	Email	string	`json:"email"`
}

func(cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request,) {
	var emailReq emailRequest

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&emailReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	userInfo, err := cfg.db.CreateUser(context.Background(), emailReq.Email)
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