package test

//
//import (
//	"encoding/json"
//	"github.com/gofiber/fiber/v2"
//	"github.com/stretchr/testify/assert"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"say-it"
//	"strings"
//	"testing"
//)
//
////var app = fiber.New()
//
//func TestRoutingHelloWorld(t *testing.T) {
//
//	main.app.Get("/", func(ctx *fiber.Ctx) error {
//		return ctx.SendString("Hello World")
//	})
//
//	request := httptest.NewRequest("GET", "/", nil)
//	response, err := main.app.Test(request)
//	assert.Nil(t, err)
//	assert.Equal(t, 200, response.StatusCode)
//
//	body, err := io.ReadAll(response.Body)
//	assert.Nil(t, err)
//	assert.Equal(t, "Hello World", string(body))
//}
//
//func TestCtx(t *testing.T) {
//
//	main.app.Get("/hello", func(ctx *fiber.Ctx) error {
//		name := ctx.Query("name", "guest")
//		return ctx.SendString("Hello " + name)
//	})
//
//	request := httptest.NewRequest("GET", "/hello?name=imam", nil)
//	response, err := main.app.Test(request)
//	assert.Nil(t, err)
//	assert.Equal(t, 200, response.StatusCode)
//
//	body, err := io.ReadAll(response.Body)
//	assert.Nil(t, err)
//	assert.Equal(t, "Hello imam", string(body))
//
//	request = httptest.NewRequest("GET", "/hello", nil)
//	response, err = main.app.Test(request)
//	assert.Nil(t, err)
//	assert.Equal(t, 200, response.StatusCode)
//
//	body, err = io.ReadAll(response.Body)
//	assert.Nil(t, err)
//	assert.Equal(t, "Hello guest", string(body))
//}
//
//func TestHttpRequest(t *testing.T) {
//
//	main.app.Get("/request", func(ctx *fiber.Ctx) error {
//		first := ctx.Get("firstname")
//		last := ctx.Cookies("lastname")
//		return ctx.SendString("Hello " + first + " " + last)
//	})
//
//	request := httptest.NewRequest("GET", "/request", nil)
//	request.Header.Set("firstname", "imam")
//	request.AddCookie(&http.Cookie{Name: "lastname", Value: "ahmad"})
//	response, err := main.app.Test(request)
//	assert.Nil(t, err)
//	assert.Equal(t, 200, response.StatusCode)
//
//	body, err := io.ReadAll(response.Body)
//	assert.Nil(t, err)
//	assert.Equal(t, "Hello imam ahmad", string(body))
//}
//
//func TestParamRequest(t *testing.T) {
//
//	main.app.Get("/users/:userId/orders/:orderId", func(ctx *fiber.Ctx) error {
//		userId := ctx.Params("userId")
//		orderId := ctx.Params("orderId")
//		return ctx.SendString("Get Order " + orderId + " From User " + userId)
//	})
//
//	request := httptest.NewRequest("GET", "/users/imam/orders/10", nil)
//	response, err := main.app.Test(request)
//	assert.Nil(t, err)
//	assert.Equal(t, 200, response.StatusCode)
//
//	body, err := io.ReadAll(response.Body)
//	assert.Nil(t, err)
//	assert.Equal(t, "Get Order 10 From User imam", string(body))
//}
//
//func TestFormRequest(t *testing.T) {
//
//	main.app.Post("/hello", func(ctx *fiber.Ctx) error {
//		name := ctx.FormValue("name")
//		return ctx.SendString("Hello " + name)
//	})
//
//	reader := strings.NewReader("name=imam")
//
//	request := httptest.NewRequest("POST", "/hello", reader)
//	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
//	response, err := main.app.Test(request)
//	assert.Nil(t, err)
//	assert.Equal(t, 200, response.StatusCode)
//
//	body, err := io.ReadAll(response.Body)
//	assert.Nil(t, err)
//	assert.Equal(t, "Hello imam", string(body))
//}
//
//type LoginRequest struct {
//	Username string `json:"username"`
//	Password string `json:"password"`
//}
//
//func TestBodyRequest(t *testing.T) {
//
//	main.app.Post("/login", func(ctx *fiber.Ctx) error {
//		body := ctx.Body()
//
//		request := new(LoginRequest)
//		err := json.Unmarshal(body, request)
//		if err != nil {
//			return err
//		}
//		return ctx.SendString("Hello " + request.Username)
//	})
//
//	reader := strings.NewReader(
//		`{"username": "imam",
//			"password": "imammo"}`,
//	)
//
//	request := httptest.NewRequest("POST", "/login", reader)
//	request.Header.Set("Content-Type", "application/json")
//
//	response, err := main.app.Test(request)
//	assert.Nil(t, err)
//	assert.Equal(t, 200, response.StatusCode)
//
//	body, err := io.ReadAll(response.Body)
//	assert.Nil(t, err)
//	assert.Equal(t, "Hello imam", string(body))
//}
//
//func TestResponseJSON(t *testing.T) {
//
//	main.app.Get("/user", func(ctx *fiber.Ctx) error {
//
//		return ctx.JSON(fiber.Map{
//			"name":     "imam ahmad",
//			"username": "imam",
//		})
//	})
//
//	request := httptest.NewRequest("GET", "/user", nil)
//	request.Header.Set("Accept", "application/json")
//
//	response, err := main.app.Test(request)
//	assert.Nil(t, err)
//	assert.Equal(t, 200, response.StatusCode)
//
//	body, err := io.ReadAll(response.Body)
//	assert.Nil(t, err)
//	assert.Equal(t, `{"name":"imam ahmad","username":"imam"}`, string(body))
//}
