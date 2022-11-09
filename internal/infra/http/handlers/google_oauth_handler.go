package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"trainee/config"
	"trainee/internal/app"
	"trainee/internal/infra/http/requests"
	"trainee/internal/infra/http/response"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type OauthHandler struct {
	us app.UserService
	as app.AuthService
}

func NewOauthHandler(u app.UserService, a app.AuthService) OauthHandler {
	return OauthHandler{
		us: u,
		as: a,
	}
}

func (o OauthHandler) GetInfo(ctx echo.Context) error {
	googleConfig := config.LoadOAUTHConfiguration()
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Google oauth error: %s", err))
	}
	state := base64.URLEncoding.EncodeToString(b)
	url := googleConfig.AuthCodeURL(state)
	log.Println(url)
	err = ctx.Redirect(http.StatusTemporaryRedirect, url)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Google oauth error, could not redirect url: %s", err))
	}
	return response.MessageResponse(ctx, http.StatusOK, "Success")
}

func (o OauthHandler) CallBackRegister(ctx echo.Context) error {
	cfg := config.LoadOAUTHConfiguration()

	token, err := cfg.Exchange(context.Background(), ctx.FormValue("code"))
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("Google oauth error, code exchange wrong: %s", err))
	}
	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("Google oauth error, failed gatting user info: %s", err))
	}
	defer resp.Body.Close()
	var usr requests.RegisterOauth2
	err = json.NewDecoder(resp.Body).Decode(&usr)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not decode user data")
	}
	//todo change password creation
	usr.Password = usr.ID + usr.Email
	userFromRegister := usr.RegisterToUser()
	user, err := o.as.Register(userFromRegister)
	if err != nil {
		if strings.HasSuffix(err.Error(), "invalid credentials user exist") {
			user, err = o.us.FindByEmail(userFromRegister.Email)
			log.Println(userFromRegister)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("User not exist: %s", err))
			}
			if user.ID == 0 || (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userFromRegister.Password)) != nil) {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
			}
		}
	}
	accessToken, exp, err := o.as.CreateAccessToken(user)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Google oauth error, coukd not create access token: %s", err))
	}
	refreshToken, err := o.as.CreateRefreshToken(user)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Google oauth error, coukd not create access token: %s", err))
	}
	res := response.NewLoginResponse(accessToken, refreshToken, exp)
	return response.Response(ctx, http.StatusOK, res)
}
