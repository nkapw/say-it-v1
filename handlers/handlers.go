package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
	"say-it/connection"
	"say-it/helper"
	"say-it/models"
	"strconv"
	"strings"
	"time"
)

var db = connection.GetConnection()

var secretKey = []byte("your_secret_key")
var validate *validator.Validate = validator.New()

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
	//router.Handle("/user/{id}", authMiddleware(http.HandlerFunc(GetUserHandler))).Methods("GET")
	//router.Handle("/user/{id}", authMiddleware(http.HandlerFunc(EditUserHandler))).Methods("PUT")
	router.Handle("/user/update", authMiddleware(http.HandlerFunc(UpdateCurrentUserHandler))).Methods("PUT")

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
	err = validate.Struct(user)

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

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var user models.User
	err := db.QueryRow("SELECT id, name, email, password, username FROM users WHERE id=$1", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Username)
	if err != nil {
		//http.Error(w, err.Error(), 404)
		response := models.NewErrorResponse("user not found", "not found")
		helper.WriteToResponseBody(w, http.StatusNotFound, &response)
		return
	}

	response := models.NewSuccessResponse("ok", user)
	helper.WriteToResponseBody(w, http.StatusOK, response)
}

func EditUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var updatedUser models.User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		response := models.NewErrorResponse("invalid request payload", "bad request")
		helper.WriteToResponseBody(w, http.StatusBadRequest, &response)
		return
	}

	_, err = db.Exec("UPDATE users SET username=$1 WHERE id=$2",
		updatedUser.Username, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.NewSuccessResponse("ok", nil)
	helper.WriteToResponseBody(w, http.StatusOK, response)
}

//func UpdateCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
//
//	// Mendapatkan ID pengguna dari token
//	userID, err := getUserIDFromToken(r)
//
//	if err != nil {
//		http.Error(w, "Invalid token", http.StatusUnauthorized)
//		return
//	}
//
//	// Mendapatkan data pengguna dari database
//	var currentUser models.User
//	err = db.QueryRow("SELECT id, name, email, password, username FROM users WHERE id=$1", userID).
//		Scan(&currentUser.ID, &currentUser.Name, &currentUser.Email, &currentUser.Password, &currentUser.Username)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusNotFound)
//		return
//	}
//
//	// Mendapatkan data yang diperbarui dari body request
//	var updatedUser models.User
//	err = json.NewDecoder(r.Body).Decode(&updatedUser)
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	// Perbarui informasi pengguna saat ini
//	_, err = db.Exec("UPDATE users SET name=$1, email=$2, password=$3, username=$4 WHERE id=$5",
//		currentUser.Name, currentUser.Email, currentUser.Password, updatedUser.Username, userID)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//}

// Fungsi helper untuk mendapatkan ID pengguna dari token
func getUserIDFromToken(r *http.Request) (int, error) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("missing Authorization header")
	}

	token := strings.Split(authHeader, " ")
	if len(token) != 2 || token[0] != "Bearer" {
		return 0, fmt.Errorf("invalid token format")
	}

	claims := jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(token[1], &claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return 0, fmt.Errorf("invalid token")
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID in token")
	}

	return userID, nil
}

// handlers/handlers.go

// handlers/handlers.go
func UpdateCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan ID pengguna dari token
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Mendapatkan data pengguna dari database
	var currentUser models.User
	var currentUserPfp sql.NullString

	if currentUserPfp.Valid {
		currentUser.ProfilePicture = currentUserPfp.String
	}
	err = db.QueryRow("SELECT id, name, email, password, username, profile_picture FROM users WHERE id=$1", userID).
		Scan(&currentUser.ID, &currentUser.Name, &currentUser.Email, &currentUser.Password, &currentUser.Username, &currentUserPfp)
	if err != nil {
		log.Fatal("hello QueryRow")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Mendapatkan data yang diperbarui dari formulir multipart
	err = r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		log.Fatal("hello ParseMultipartForm")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mengambil file gambar dari formulir
	file, header, err := r.FormFile("profile_picture")
	if err == nil {
		// Jika ada file gambar, simpan di Google Cloud Storage
		defer file.Close()

		// Inisialisasi klien GCS
		gcsClient, err := helper.CreateGCSClient()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gcsClient.Close()

		// Simpan gambar di GCS
		bucketName := "profile_picture_bucket"
		objectName := fmt.Sprintf("profile_%d_%s", userID, header.Filename)

		ctx := context.Background()
		wc := gcsClient.Bucket(bucketName).Object(objectName).NewWriter(ctx)
		_, err = io.Copy(wc, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := wc.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Dapatkan URL gambar GCS
		imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

		// Simpan URL gambar di database
		_, err = db.Exec("UPDATE users SET profile_picture=$1 WHERE id=$2", imageURL, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update informasi pengguna
		currentUser.ProfilePicture = imageURL
		response := models.NewSuccessResponse("User information updated successfully", currentUser)
		helper.WriteToResponseBody(w, http.StatusOK, response)
	}
	//
	//// Mendapatkan data pengguna yang diperbarui dari body request
	//var updatedUser models.User
	//err = json.NewDecoder(r.Body).Decode(&updatedUser)
	//if err != nil {
	//	//log.Fatal(err.Error())
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//
	//// Perbarui informasi pengguna saat ini
	//_, err = db.Exec("UPDATE users SET name=$1, email=$2, password=$3, username=$4 WHERE id=$5",
	//	updatedUser.Name, updatedUser.Email, updatedUser.Password, updatedUser.Username, userID)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	//// Update informasi pengguna di respons
	//currentUser.Name = updatedUser.Name
	//currentUser.Email = updatedUser.Email
	//currentUser.Username = updatedUser.Username
	//log.Fatal("selesai")
	//
	//// Kirim respons JSON
	//response := models.NewSuccessResponse("User information updated successfully", currentUser)
	//helper.WriteToResponseBody(w, http.StatusOK, response)
}
