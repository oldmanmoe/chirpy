package main

import (
	internal "chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

	var day = time.Hour * 24
	var accessExpireTime = time.Hour
	var refreshExpireTime = day * 60
	var now = time.Now()
type User struct {
	ID        		uuid.UUID 	`json:"id"`
	CreatedAt 		time.Time 	`json:"created_at"`
	UpdatedAt 		time.Time 	`json:"updated_at"`
	Email     		string    	`json:"email"`
	Token			string		`json:"token"`
	RefreshToken 	string		`json:"refresh_token"`
}

type AuthRequest struct {
	Email            string  	`json:"email"`
	Password         string  	`json:"password"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var registerReq AuthRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&registerReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	password, err := internal.HashPassword(registerReq.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong hashing password")
	}



	userInfo, err := cfg.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			Email:          registerReq.Email,
			HashedPassword: password,
		})
	if err != nil {
		log.Fatal(err)
		return
	}

	newUser := User{
		ID:        userInfo.ID,
		CreatedAt: userInfo.CreatedAt,
		UpdatedAt: userInfo.UpdatedAt,
		Email:     userInfo.Email,
	}

	printUserInfo(newUser)

	respondWithJSON(w, http.StatusCreated, newUser)
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	

	var loginReq AuthRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
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

	accessToken, _ := internal.MakeJWT(existingUser.ID, cfg.secret, accessExpireTime) 
	refreshToken := internal.MakeRefreshToken()

	err = cfg.db.StoreNewRefreshTokenInfo(
		context.Background(),
		database.StoreNewRefreshTokenInfoParams{
			Token: refreshToken,
			UserID: existingUser.ID,
			ExpiresAt: now.Add(refreshExpireTime),
		},
	)

	if err != nil {
		fmt.Printf("An error has occurred: %v", err)
	}

	resultUser := User{
		ID:        		existingUser.ID,
		CreatedAt: 		existingUser.CreatedAt,
		UpdatedAt: 		existingUser.UpdatedAt,
		Email:     		existingUser.Email,
		Token:	   		accessToken,
		RefreshToken:	refreshToken,
	}

	respondWithJSON(w, http.StatusOK, resultUser)
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	
	reqRefreshToken, err := internal.GetBearerToken(r.Header)
	if err != nil {
		fmt.Print(err)
		return 
	}

	dbRefreshToken, err := cfg.db.GetUserFromRefreshToken(context.Background(), reqRefreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized, authentication required")
		return
	}

	if dbRefreshToken.RevokedAt.Valid == true {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized, authentication required.")
	}

	if now.After(dbRefreshToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Session expired")
		return
	}

	newAccessToken, err := internal.MakeJWT(dbRefreshToken.UserID, cfg.secret, accessExpireTime)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized, authentication required")
	}
	
	result := struct {
		Token	string	`json:"token"`
	}{
		Token: newAccessToken,
	}

	respondWithJSON(w, http.StatusOK, result)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

	reqRefreshToken, err := internal.GetBearerToken(r.Header)
	if err != nil {
		fmt.Print(err)
		return 
	}
	
	dbRefreshToken, err := cfg.db.GetUserFromRefreshToken(context.Background(), reqRefreshToken)
	if err != nil {
		fmt.Print(err)
		return 
	}

	cfg.db.UpdateRevokedAt(
		context.Background(),
		database.UpdateRevokedAtParams{
				RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
				UpdatedAt: time.Now(),
				Token: dbRefreshToken.Token,
			},
		)
	

	respondWithJSON(w, http.StatusNoContent, "")

}

func (cfg *apiConfig) handlerUpdatePassword(w http.ResponseWriter, r *http.Request) {

	var reqUserInfo AuthRequest
	reqRefreshToken, err := internal.GetBearerToken(r.Header)
	if err != nil {
		fmt.Print(err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized, authentication required")
		return
	}

	reqUser, err := internal.ValidateJWT(reqRefreshToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized, authentication required")
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqUserInfo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	newPassword, err := internal.HashPassword(reqUserInfo.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	userInfoUpdate, err := cfg.db.UpdateUserEmailAndPassword(
		context.Background(),
		database.UpdateUserEmailAndPasswordParams{
			Email: 			reqUserInfo.Email,
			HashedPassword: newPassword,
			ID:				reqUser,
		})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Something went wrong")
		return
	}

	updatedUserInfo := User{
		ID: 		userInfoUpdate.ID,
		CreatedAt: 	userInfoUpdate.CreatedAt,
		UpdatedAt: 	time.Now(),
		Email: 		userInfoUpdate.Email,
	}

	printUserInfo(updatedUserInfo)

	respondWithJSON(w, http.StatusOK, updatedUserInfo)

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
