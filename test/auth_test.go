package test

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
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
