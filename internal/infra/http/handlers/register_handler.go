package handlers

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"trainee/internal/app"
	"trainee/internal/infra/http/requests"
	"trainee/internal/infra/http/response"
)

type RegisterHandler struct {
	us app.UserService
	as app.AuthService
}

func NewRegisterHandler(u app.UserService, a app.AuthService) RegisterHandler {
	return RegisterHandler{
		us: u,
		as: a,
	}
}

func (r RegisterHandler) Register(c echo.Context) error {
	var registerUser requests.RegisterAuth
	if err := c.Bind(&registerUser); err != nil {
		return err
	}
	if err := c.Validate(&registerUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Empty or not valid")
	}

	userFromRegister := registerUser.RegisterToUser()

	user, err := r.as.Register(userFromRegister)
	if err != nil {
		log.Printf("Register Handler: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, user)
}

func (r RegisterHandler) Login(c echo.Context) error {
	var authUser requests.LoginAuth
	if err := c.Bind(&authUser); err != nil {
		return err
	}
	if err := c.Validate(&authUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "files not empty")
	}
	user, err := r.us.FindByEmail(authUser.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user not exist")
	}
	if user.ID == 0 || (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authUser.Password)) != nil) {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}
	accessToken, exp, err := r.as.CreateAccessToken(user)
	if err != nil {
		return err
	}
	refreshToken, err := r.as.CreateRefreshToken(user)
	if err != nil {
		return err
	}
	res := response.NewLoginResponse(accessToken, refreshToken, exp)

	return c.JSON(http.StatusOK, res)
}
