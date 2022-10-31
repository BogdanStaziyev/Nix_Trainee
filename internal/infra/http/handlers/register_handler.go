package handlers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"trainee/config"
	"trainee/internal/app"
	"trainee/internal/domain"
	"trainee/internal/infra/responses"
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
	var registerUser domain.User
	if err := c.Bind(&registerUser); err != nil {
		return err
	}
	if err := c.Validate(&registerUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Empty or not valid")
	}
	user, token, err := r.as.Register(registerUser)
	if err != nil {
		log.Printf("Register Handler: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	type userDto struct {
		User  domain.User
		Token string
	}
	ne := userDto{
		User:  user,
		Token: token,
	}
	return c.JSON(http.StatusCreated, ne)
}

func (r RegisterHandler) Login(c echo.Context) error {
	var authUser domain.AuthUser
	if err := c.Bind(&authUser); err != nil {
		return err
	}
	if err := c.Validate(&authUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "files not empty")
	}
	user, err := r.us.FindByEmail(authUser.Email)
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
	res := responses.NewLoginResponse(accessToken, refreshToken, exp)

	return c.JSON(http.StatusOK, res)
}

func (r RegisterHandler) RefreshToken(c echo.Context) error {
	request := new(domain.Refresh)
	if err := c.Bind(request); err != nil {
		return err
	}

	token, err := jwt.Parse(request.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfiguration().AccessSecret), nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}
	user, err := r.us.FindByID(int64(claims["id"].(float64)))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
	}

	accessToken, exp, err := r.as.CreateAccessToken(user)
	if err != nil {
		return err
	}
	refreshToken, err := r.as.CreateRefreshToken(user)
	if err != nil {
		return err
	}
	res := responses.NewLoginResponse(accessToken, refreshToken, exp)

	return c.JSON(http.StatusOK, res)
}
