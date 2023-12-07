package helper

import "github.com/go-playground/validator/v10"

var SecretKey = []byte("your_secret_key")
var Validate *validator.Validate = validator.New()
