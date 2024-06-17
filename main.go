package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/itworldwithrahul/golang-app/config"
	"github.com/labstack/echo"
)

type Request struct {
	Name      string     `json:"name" `
	Email     string     `json:"email" `
	Adharcard string     `json:"adharcard" `
	Age       uint8      `json:"age" validate:"gte=1,lte=130"`
	Gender    string     ` json:"gender" validate:"oneof=male female prefer_not_to"`
	Addresses []*Address `json:"address" validate:"required,dive,required"`
}

// Address houses a users address information
type Address struct {
	Street  string `json:"street" validate:"required"`
	City    string `json:"city" validate:"required"`
	Country string `json:"country" validate:"required"`
	Phone   string `json:"phone" validate:"required,len=10"`
}

type Response struct {
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Adharcard string     `json:"adharcard"`
	Age       uint8      `json:"age"`
	Gender    string     ` json:"gender"`
	Addresses []*Address `json:"address"`
	CreatedBy string     `json:"createdBy"`
	UpdatedBy string     `json:"updatedBy"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
type User struct {
	Id       string `json:"Id" validate:"required"`
	Username string `json:"UserName" validate:"required"`
	Password string `json:"Password" validate:"required"`
}
type Login struct {
	UserName string `json:"UserName" validate:"required"`
	Password string `json:"Password" validate:"required"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// CreateUser
func CreateUser(c echo.Context) error {

	//u

	// convert json to struct [ Binding/unmarshling]
	req := new(Request)
	if err := c.Bind(req); err != nil {
		log.Printf("failed to process user request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// validate request fields
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var res Response
	res.Adharcard = req.Adharcard
	res.Name = req.Name
	res.Email = req.Email

	// Not Recommended: response will point to same slice
	// res.Addresses = req.Addresses
	// log.Printf("Addreess of res.Addresses: %p\n", res.Addresses)
	// log.Printf("Addreess of req.Addresses: %p\n", req.Addresses)

	res.Age = req.Age
	res.Gender = req.Gender
	res.Addresses = append(res.Addresses, req.Addresses...) // wont'work
	// log.Printf("Addreess of res.Addresses: %p\n", res.Addresses)
	// log.Printf("Addreess of req.Addresses: %p\n", req.Addresses)

	res.CreatedAt = time.Now()
	res.UpdatedAt = time.Now()
	//ToDo: parse info from token then initialize
	res.CreatedBy = "user"
	res.UpdatedBy = "user"

	return c.JSON(http.StatusCreated, res)
}

// HelloServer
func HelloServer(c echo.Context) error {
	return c.String(http.StatusOK, "Hello server")
}

// GetUser
func GetUser(c echo.Context) error {
	userid := c.Param("userid")
	return c.JSON(http.StatusOK, map[string]string{
		"id":      userid,
		"message": "Welcome user",
	})
}

func loginUser(c echo.Context) error {

	req := new(Login)
	if err := c.Bind(req); err != nil {
		log.Printf("failed to process user request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// validate request fields
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var res Login
	res.UserName = req.UserName
	res.Password = req.Password

	return c.JSON(http.StatusCreated, res)
}

func DeleteUser(c echo.Context) error {
	userid := c.Param("userid")
	return c.JSON(http.StatusOK, map[string]string{
		"id":      userid,
		"message": "User deleted",
	})
}

func updatedUser(c echo.Context) error {
	userid := c.Param("userid")
	log.Println("Userid", userid)

	// convert json to struct [ Binding/unmarshling]
	req := new(Request)
	if err := c.Bind(req); err != nil {
		log.Printf("failed to process user request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// validate request fields
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var res Response
	res.Adharcard = req.Adharcard
	res.Name = req.Name
	res.Email = req.Email

	// Not Recommended: response will point to same slice
	// res.Addresses = req.Addresses
	// log.Printf("Addreess of res.Addresses: %p\n", res.Addresses)
	// log.Printf("Addreess of req.Addresses: %p\n", req.Addresses)

	res.Age = req.Age
	res.Gender = req.Gender
	res.Addresses = append(res.Addresses, req.Addresses...) // wont'work
	// log.Printf("Addreess of res.Addresses: %p\n", res.Addresses)
	// log.Printf("Addreess of req.Addresses: %p\n", req.Addresses)

	res.UpdatedAt = time.Now()
	//ToDo: parse info from token then initialize
	res.UpdatedBy = "user"

	return c.JSON(http.StatusCreated, res)
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log.Printf("configuration:%+v\n", cfg)

	e.GET("/", HelloServer)                // http://127.0.0.1
	e.GET("/api/v1/user/:userid", GetUser) // http://127.0.0.1/api/v1/user/123
	e.POST("/api/v1/user", CreateUser)     // http://127.0.0.1/api/v1/user/123
	e.POST("/api/v1/user/login", loginUser)
	e.PUT("/api/v1/user/:id", updatedUser)
	e.DELETE("/api/v1/user/:id", DeleteUser)
	// start the webserver
	e.Start(cfg.Server.Port)
}
