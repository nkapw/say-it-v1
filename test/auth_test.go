package test

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"say-it/handlers"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {

	bodyReq := strings.NewReader(`{"username":"imam", "name":"imam ahmad", "email":"imam@gmail.com", "password":"imampw"'`)
	request := httptest.NewRequest("POST", "/register", bodyReq)

	request.Header.Set("Content-Type", "application/json")

	//res, err := main.App.Test(request)
	res, err := fiber.New().Test(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	//
	//body, err := io.ReadAll(res.Body)
	//assert.Nil(t, err)
	//assert.Equal(t, "Hello imam", string(body))

}

func TestUpdate(t *testing.T) {

	request := httptest.NewRequest("POST", "/user/update", nil)

	request.Header.Set("Au", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE3OTMwMjcsImlhdCI6MTcwMTc4OTQyNywic3ViIjoiMiJ9.sjtuQ2Tj7RZ4U06DZ5sTL3UzuLp-WZfy7kJn1GsvAsA")

	recorder := httptest.NewRecorder()

	handlers.UpdateCurrentUserHandler(recorder, request)

	response := recorder.Result()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))

}
