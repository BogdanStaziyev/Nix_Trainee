package handlers

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"trainee/internal/app"
	"trainee/internal/domain"
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

// Register 		godoc
// @Summary 		Register
// @Description 	New user registration
// @ID				user-register
// @Tags			Auth Actions
// @Accept 			json
// @Produce 		json
// @Param			input body requests.RegisterAuth true "users email, users password"
// @Success 		201 {object} response.UserResponse
// @Failure			400 {object} response.Error
// @Router			/register [post]
func (r RegisterHandler) Register(ctx echo.Context) error {
	var registerUser requests.RegisterAuth
	if err := ctx.Bind(&registerUser); err != nil {
		return err
	}
	if err := ctx.Validate(&registerUser); err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Empty or not valid")
	}

	userFromRegister := registerUser.RegisterToUser()

	user, err := r.as.Register(userFromRegister)
	if err != nil {
		log.Printf("Register Handler: %s", err)
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	userResponse := domain.User.DomainToResponse(user)
	return response.Response(ctx, http.StatusCreated, userResponse)
}

// Login 			godoc
// @Summary 		LoginAuth
// @Description 	LoginAuth
// @Tags			Auth Actions
// @Accept 			json
// @Produce 		json
// @Param			input body requests.LoginAuth true "users email, users password"
// @Success 		201 {object} response.LoginResponse
// @Failure			400 {object} response.Error
// @Failure			401 {object} response.Error
// @Router			/login [post]
func (r RegisterHandler) Login(ctx echo.Context) error {
	var authUser requests.LoginAuth
	if err := ctx.Bind(&authUser); err != nil {
		return err
	}
	if err := ctx.Validate(&authUser); err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "files not empty")
	}
	user, err := r.us.FindByEmail(authUser.Email)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "user not exist")
	}
	if user.ID == 0 || (bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authUser.Password)) != nil) {
		return response.ErrorResponse(ctx, http.StatusUnauthorized, "invalid credentials")
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

	return response.Response(ctx, http.StatusOK, res)
}
