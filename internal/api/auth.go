package api

import (
	"encoding/json"
	"net/http"
	"packhaus/internal/db"
	"packhaus/internal/utils"
	"strconv"
)

type signupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signupResponse struct {
	Token string  `json:"token"`
	User  db.User `json:"user"`
}

type signinRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signinResponse struct {
	Token string `json:"token"`
}

func (cntlr *controller) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user, err := db.CreateUser(cntlr.DB, req.Username, req.Email, string(hashed))
	if err != nil {
		http.Error(w, "user creation failed", http.StatusInternalServerError)
		return
	}

	token, err := utils.GenerateJWT(strconv.Itoa(user.ID))
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(signupResponse{
		Token: token,
		User:  user,
	})
}

func (cntlr *controller) SigninHandler(w http.ResponseWriter, r *http.Request) {
	var req signinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByUsername(cntlr.DB, req.Username)
	if err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	if err := utils.VerifyPassword(user.Password, req.Password); err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(strconv.Itoa(user.ID))
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(signinResponse{
		Token: token,
	})
}
