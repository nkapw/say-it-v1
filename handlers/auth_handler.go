package handlers

//
//import (
//	"encoding/json"
//	"github.com/golang-jwt/jwt/v4"
//	"golang.org/x/crypto/bcrypt"
//	"net/http"
//	"say-it/connection"
//	"say-it/models"
//	"strconv"
//	"time"
//)
//
//var db = connection.GetConnection()
//
//var secretKey = []byte("your_secret_key")
//
//func RegisterHandler(w http.ResponseWriter, r *http.Request) {
//	w.Header().Add("Content-Type", "application/json")
//	var user models.User
//	err := json.NewDecoder(r.Body).Decode(&user)
//	if err != nil {
//		//http.Error(w, err.Error(), http.StatusBadRequest)
//
//		w.WriteHeader(http.StatusBadRequest)
//
//		res := map[string]string{
//			"status": "bad request",
//			//"message": "email or username already exist",
//		}
//		json.NewEncoder(w).Encode(&res)
//		return
//	}
//
//	rows, err := db.Query("SELECT username, email FROM users WHERE username = $1 OR email = $2", user.Username, user.Email)
//	if err != nil {
//		//http.Error(w, err.Error(), http.StatusInternalServerError)
//		w.WriteHeader(http.StatusInternalServerError)
//
//		res := map[string]string{
//			"status": "internal server error",
//			//"message": "email or username already exist",
//		}
//		json.NewEncoder(w).Encode(&res)
//		return
//	}
//	var username, email string
//	if rows.Next() {
//		err := rows.Scan(&username, &email)
//		if err != nil {
//			//http.Error(w, err.Error(), http.StatusInternalServerError)
//			w.WriteHeader(http.StatusInternalServerError)
//
//			res := map[string]string{
//				"status": "internal server error",
//				//"message": "email or username already exist",
//			}
//			json.NewEncoder(w).Encode(&res)
//			return
//		}
//	}
//
//	if username == user.Username || email == user.Email {
//		//http.Error(w, "email or username already exist", http.StatusBadRequest)
//		//response := map[string]string{"token": "message":"email or username already exist"}
//		w.WriteHeader(http.StatusBadRequest)
//
//		res := map[string]string{
//			"status":  "bad request",
//			"message": "email or username already exist",
//		}
//		json.NewEncoder(w).Encode(&res)
//
//		return
//	}
//
//	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//
//		res := map[string]string{
//			"status": "internal server error",
//			//"message": "email or username already exist",
//		}
//		json.NewEncoder(w).Encode(&res)
//		return
//	}
//	user.Password = string(hashedPassword)
//
//	_, err = db.Exec("INSERT INTO users (name, email, password, username) VALUES ($1, $2, $3, $4)",
//		user.Name, user.Email, user.Password, user.Username)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//
//		res := map[string]string{
//			"status": "internal server error",
//			//"message": "email or username already exist",
//		}
//		json.NewEncoder(w).Encode(&res)
//		return
//	}
//
//	w.WriteHeader(http.StatusCreated)
//}
//
//func LoginHandler(w http.ResponseWriter, r *http.Request) {
//	w.Header().Add("Content-Type", "application/json")
//	var loginReq models.User
//	err := json.NewDecoder(r.Body).Decode(&loginReq)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//
//		res := map[string]string{
//			"status": "bad request",
//			//"message": "email or username already exist",
//		}
//		json.NewEncoder(w).Encode(&res)
//
//		return
//	}
//
//	var dbUser models.User
//	err = db.QueryRow("SELECT id, name, email, password, username FROM users WHERE email=$1", loginReq.Email).
//		Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email, &dbUser.Password, &dbUser.Username)
//	if err != nil {
//		//http.Error(w, "Invalid email or password", http.StatusUnauthorized)
//		w.WriteHeader(http.StatusUnauthorized)
//
//		res := map[string]string{
//			"status":  "unauthorize",
//			"message": "Invalid email or password",
//		}
//		json.NewEncoder(w).Encode(&res)
//		return
//	}
//
//	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginReq.Password))
//	if err != nil {
//		//http.Error(w, "Invalid email or password", http.StatusUnauthorized)
//		w.WriteHeader(http.StatusUnauthorized)
//
//		res := map[string]string{
//			"status":  "unauthorize",
//			"message": "Invalid email or password",
//		}
//		json.NewEncoder(w).Encode(&res)
//		return
//	}
//
//	token, err := generateToken(dbUser.ID)
//	if err != nil {
//		//http.Error(w, err.Error(), http.StatusInternalServerError)
//		w.WriteHeader(http.StatusInternalServerError)
//
//		res := map[string]string{
//			"status": "internal server error",
//			//"message": "email or username already exist",
//		}
//		json.NewEncoder(w).Encode(&res)
//		return
//	}
//
//	response := map[string]string{"token": token}
//	json.NewEncoder(w).Encode(response)
//}
//
//func generateToken(userID int) (string, error) {
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
//		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
//		IssuedAt:  time.Now().Unix(),
//		Subject:   strconv.Itoa(userID),
//	})
//
//	signedToken, err := token.SignedString(secretKey)
//	if err != nil {
//		return "", err
//	}
//
//	return signedToken, nil
//}
