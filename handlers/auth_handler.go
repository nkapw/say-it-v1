package handlers

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"say-it/helper"
	"say-it/models"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	err = helper.Validate.Struct(user)

	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)

		response := models.NewErrorResponse("invalid request payload", "bad request")
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
		response := models.NewErrorResponse("email or username already exist", "bad request")
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

	_, err = db.Exec("INSERT INTO users (name, email, password, username) VALUES ($1, $2, $3, $4)",
		user.Name, user.Email, user.Password, user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var registerResponse = struct {
		Name  string
		Email string
	}{
		Name:  user.Name,
		Email: user.Email,
	}

	response := models.NewSuccessResponse("success registered", registerResponse)
	helper.WriteToResponseBody(w, http.StatusCreated, &response)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var loginReq models.User
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		response := models.NewErrorResponse("invalid request payload", "bad request")
		helper.WriteToResponseBody(w, http.StatusBadRequest, &response)
		return
	}

	var dbUser models.User
	err = db.QueryRow("SELECT id, name, email, password, username FROM users WHERE email=$1", loginReq.Email).
		Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email, &dbUser.Password, &dbUser.Username)
	if err != nil {
		response := models.NewErrorResponse("Invalid email or password", "unauthorized")
		helper.WriteToResponseBody(w, http.StatusUnauthorized, &response)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginReq.Password))
	if err != nil { //http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		response := models.NewErrorResponse("Invalid email or password", "unauthorized")
		helper.WriteToResponseBody(w, http.StatusUnauthorized, &response)
		return
	}

	token, err := generateToken(dbUser.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//response := map[string]string{"token": token}
	//json.NewEncoder(w).Encode(response)

	response := models.NewSuccessResponse("ok", token)
	helper.WriteToResponseBody(w, http.StatusOK, response)
}
