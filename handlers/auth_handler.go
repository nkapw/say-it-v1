package handlers

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"say-it/helper"
	"say-it/models"
	"strconv"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	err = helper.Validate.Struct(user)
	if err != nil {
		response := models.NewErrorResponse("Invalid Request Payload", "bad request", err.Error())
		helper.WriteToResponseBody(w, http.StatusBadRequest, &response)
		return
	}

	rows, err := db.Query("SELECT username, email FROM users WHERE username = $1 OR email = $2", user.Username, user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var username, email string
	if rows.Next() {
		err := rows.Scan(&username, &email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if username == user.Username || email == user.Email {
		response := models.NewErrorResponse("Registration failed", "Bad Request", "Username or Email is already use")
		helper.WriteToResponseBody(w, http.StatusBadRequest, &response)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	_, err = db.Exec("INSERT INTO users ( email, password, username) VALUES ($1, $2, $3)",
		user.Email, user.Password, user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	registerResponse := models.RegisterResponse{
		Username: user.Username,
		Email:    user.Email,
	}

	response := models.NewSuccessResponse("Registration successful", registerResponse)
	helper.WriteToResponseBody(w, http.StatusCreated, &response)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var loginReq models.User
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		response := models.NewErrorResponse("Invalid Request Payload", "Bad Request", err.Error())
		helper.WriteToResponseBody(w, http.StatusBadRequest, &response)
		return
	}

	var dbUser models.User
	err = db.QueryRow("SELECT id, email, password, username FROM users WHERE email=$1", loginReq.Email).
		Scan(&dbUser.ID, &dbUser.Email, &dbUser.Password, &dbUser.Username)
	if err != nil {
		response := models.NewErrorResponse("Login Failed", "Unauthorized", "Invalid Email")
		helper.WriteToResponseBody(w, http.StatusUnauthorized, &response)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginReq.Password))
	if err != nil { //http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		response := models.NewErrorResponse("Login Failed", "Unauthorized", "Invalid Password")
		helper.WriteToResponseBody(w, http.StatusUnauthorized, &response)
		return
	}

	token, err := generateToken(dbUser.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loginResponse := models.LoginResponse{
		Id:       strconv.Itoa(dbUser.ID),
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Token:    token,
	}

	response := models.NewSuccessResponse("Login successful", loginResponse)
	helper.WriteToResponseBody(w, http.StatusOK, response)
}
