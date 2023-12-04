package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"say-it/connection"
	"say-it/models"
	"strconv"
	"strings"
	"time"
)

//type ApiResponse struct {
//	statusCode int
//	message    string
//	//data       models.User
//}

var db = connection.GetConnection()

var secretKey = []byte("your_secret_key")

func generateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.Itoa(userID),
	})

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func isValidToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi metode tanda tangan
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	return err == nil && token.Valid
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	authMiddleware := AuthMiddleware

	router.HandleFunc("/register", RegisterHandler).Methods("POST")
	router.HandleFunc("/login", LoginHandler).Methods("POST")
	router.Handle("/user/{id}", authMiddleware(http.HandlerFunc(GetUserHandler))).Methods("GET")
	router.Handle("/user/{id}", authMiddleware(http.HandlerFunc(EditUserHandler))).Methods("PUT")

	return router
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.Split(authHeader, " ")
		if len(token) != 2 || token[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		if !isValidToken(token[1]) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rows, err := db.Query("SELECT username, email FROM users WHERE username = $1 OR email = $2", user.Username, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var username, email string
	if rows.Next() {
		err := rows.Scan(&username, &email)
		if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if username == user.Username || email == user.Email {
		http.Error(w, "email or username already exist", http.StatusBadRequest)
		//response := map[string]string{"token": "message":"email or username already exist"}
		w.WriteHeader(http.StatusBadRequest)

		res := map[string]string{
			"status":  "bad request",
			"message": "email or username already exist",
		}
		json.NewEncoder(w).Encode(&res)

		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		res := map[string]string{
			"status": "internal server error",
			//"message": "email or username already exist",
		}
		json.NewEncoder(w).Encode(&res)
		return
	}
	user.Password = string(hashedPassword)

	_, err = db.Exec("INSERT INTO users (name, email, password, username) VALUES ($1, $2, $3, $4)",
		user.Name, user.Email, user.Password, user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("aaaaaaaaaaaaaa")

		res := map[string]string{
			"status": "internal server error",
			//"message": "email or username already exist",
		}
		json.NewEncoder(w).Encode(&res)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var loginReq models.User
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		res := map[string]string{
			"status": "bad request",
			//"message": "email or username already exist",
		}
		json.NewEncoder(w).Encode(&res)

		return
	}

	var dbUser models.User
	err = db.QueryRow("SELECT id, name, email, password, username FROM users WHERE email=$1", loginReq.Email).
		Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email, &dbUser.Password, &dbUser.Username)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)

		res := map[string]string{
			"status":  "unauthorize",
			"message": "Invalid email or password",
		}
		json.NewEncoder(w).Encode(&res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginReq.Password))
	if err != nil { //http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)

		res := map[string]string{
			"status":  "unauthorize",
			"message": "Invalid email or password",
		}
		json.NewEncoder(w).Encode(&res)
		return
	}

	token, err := generateToken(dbUser.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"token": token}
	json.NewEncoder(w).Encode(response)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var user models.User
	err := db.QueryRow("SELECT id, name, email, password, username FROM users WHERE id=$1", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func EditUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var updatedUser models.User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE users SET name=$1, email=$2, password=$3, username=$4 WHERE id=$5",
		updatedUser.Name, updatedUser.Email, updatedUser.Password, updatedUser.Username, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
