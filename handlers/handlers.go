package handlers

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"net/http"
	"say-it/connection"
	"say-it/helper"
	"say-it/middleware"
	"strconv"
	"time"
)

var db = connection.GetConnection()

func generateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.Itoa(userID),
	})

	signedToken, err := token.SignedString(helper.SecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	authMiddleware := middleware.AuthMiddleware

	router.Handle("/words", authMiddleware(http.HandlerFunc(GetAllWordsHandler))).Methods("GET").Queries("page", "{page}")
	router.Handle("/words/{WordID}", authMiddleware(http.HandlerFunc(GetWordDetailHandler))).Methods("GET")
	// router.HandleFunc("/words", GetAllWordsHandler).Methods("GET").Queries("page", "{page}") // for testing in postman
	// router.HandleFunc("/words/{WordID}", GetWordDetailHandler).Methods("GET") //for testing in postman
	router.HandleFunc("/register", RegisterHandler).Methods("POST")
	router.HandleFunc("/login", LoginHandler).Methods("POST")
	router.Handle("/user", authMiddleware(http.HandlerFunc(GetUserHandler))).Methods("GET")
	//router.Handle("/user/{id}", authMiddleware(http.HandlerFunc(EditUserHandler))).Methods("PUT")
	router.Handle("/user/update", authMiddleware(http.HandlerFunc(UpdateCurrentUserHandler))).Methods("PUT")
	router.Handle("/words/{WordID}", authMiddleware(http.HandlerFunc(GradingHandler))).Methods("POST")

	return router
}
