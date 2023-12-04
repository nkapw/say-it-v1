// package main
//
// import (
//
//	"github.com/gofiber/fiber/v2"
//	"golang.org/x/crypto/bcrypt"
//	"net/http"
//	"say-it/connection"
//
// )
//
//	type User struct {
//		Id       string `json:"id"`
//		Username string `json:"username"`
//		Name     string `json:"name"`
//		Email    string `json:"email"`
//		Password string `json:"password"`
//	}
//
// var app = fiber.New()
//
// func main() {
//
//		app.Post("/register", Register)
//
//		app.Listen("localhost:8080")
//	}
//
//	func Register(ctx *fiber.Ctx) error {
//		db := connection.GetConnection()
//		var user User
//		if err := ctx.BodyParser(&user); err != nil {
//			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
//		}
//
//		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
//		if err != nil {
//			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
//		}
//
//		_, err = db.Exec("INSERT INTO users (username, name, email, password) VALUES ($1, $2, $3, $4)", user.Username, user.Name, user.Email, string(hashedPassword))
//		if err != nil {
//			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
//		}
//
//		return ctx.Status(http.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
//
// }

package main

import (
	"log"
	"net/http"
	"say-it/handlers"
)

func main() {
	router := handlers.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
