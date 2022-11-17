package handlers_test

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/upper/db/v4"
	"net/http"
	"testing"
	"trainee/internal/app"
	"trainee/internal/app/mocks"
	"trainee/internal/domain"
	"trainee/internal/infra/http/handlers"
	"trainee/internal/infra/http/handlers/test_case"
	"trainee/internal/infra/http/requests"
)

func TestRegisterHandler_Register(t *testing.T) {
	userMockRequest := requests.RegisterAuth{
		Email:    "user@mail.com",
		Name:     "Name",
		Password: "qwerty1234",
	}

	userMockDomain := domain.User{
		ID:       1,
		Email:    "user@mail.com",
		Name:     "Name",
		Password: "qwerty1234",
	}

	requestRegister := test_case.Request{
		Method: http.MethodPost,
		Url:    "/register",
	}
	handleSuccessCreate := func(c echo.Context) error {
		mockAuth := func(user requests.RegisterAuth) app.AuthService {
			mock := mocks.NewAuthService(t)
			userDomain := domain.User{
				Email:    user.Email,
				Name:     user.Name,
				Password: user.Password,
			}
			mock.On("Register", userDomain).Return(userMockDomain, nil).Times(1)
			return mock
		}(userMockRequest)
		return handlers.NewRegisterHandler(mockAuth).Register(c)
	}

	handleErrorCreate := func(c echo.Context) error {
		mockAuth := func(user requests.RegisterAuth) app.AuthService {
			mock := mocks.NewAuthService(t)
			userDomain := domain.User{
				Email:    user.Email,
				Name:     user.Name,
				Password: user.Password,
			}
			mock.On("Register", userDomain).Return(domain.User{}, db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(userMockRequest)
		return handlers.NewRegisterHandler(mockAuth).Register(c)
	}

	handleMock := func(c echo.Context) error {
		mockAuth := func() app.AuthService {
			return mocks.NewAuthService(t)
		}()
		return handlers.NewRegisterHandler(mockAuth).Register(c)
	}

	cases := []test_case.TestCase{
		{
			"RegisterUser success",
			requestRegister,
			userMockRequest,
			handleSuccessCreate,
			test_case.ExpectedResponse{
				StatusCode: 201,
				BodyPart:   "{\"id\":1,\"email\":\"user@mail.com\",\"name\":\"Name\"}\n"},
		},
		{
			"RegisterUser error",
			requestRegister,
			userMockRequest,
			handleErrorCreate,
			test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not save new user: upper: collection does not exist\"}\n"},
		},
		{
			"Error decode user data",
			requestRegister,
			"",
			handleMock,
			test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not decode user data\"}\n"},
		},
		{
			"Error validate user data",
			requestRegister,
			requests.RegisterAuth{
				Email:    "user@email.com",
				Password: "qwerty123",
				Name:     "",
			},
			handleMock,
			test_case.ExpectedResponse{
				StatusCode: 422,
				BodyPart:   "{\"code\":422,\"error\":\"Could not validate user data\"}\n"},
		},
	}
	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			c, recorder := test_case.PrepareContextFromTestCase(test)
			c.Set("user", test_case.Token())

			if assert.NoError(t, test.HandlerFunc(c)) {
				assert.Contains(t, recorder.Body.String(), test.Expected.BodyPart)
				assert.Equal(t, recorder.Code, test.Expected.StatusCode)
			}
		})
	}
}

func TestRegisterHandler_Login(t *testing.T) {
	userMockRequest := requests.LoginAuth{
		Email:    "user@mail.com",
		Password: "qwerty1234",
	}

	userMockDomain := domain.User{
		ID:       1,
		Email:    "user@mail.com",
		Name:     "Name",
		Password: "qwerty1234",
	}

	requestRegister := test_case.Request{
		Method: http.MethodPost,
		Url:    "/login",
	}
	handleSuccessLogin := func(c echo.Context) error {
		mockAuth := func(user requests.LoginAuth) app.AuthService {
			mock := mocks.NewAuthService(t)
			mock.On("Login", user).Return(userMockDomain, "refresh", nil).Times(1).
				On("CreateAccessToken", userMockDomain).Return("access", int64(123), nil)
			return mock
		}(userMockRequest)
		return handlers.NewRegisterHandler(mockAuth).Login(c)
	}

	handleErrorLoginNoMoreRows := func(c echo.Context) error {
		mockAuth := func(user requests.LoginAuth) app.AuthService {
			mock := mocks.NewAuthService(t)
			mock.On("Login", user).Return(domain.User{}, "", db.ErrNoMoreRows).Times(1)
			return mock
		}(userMockRequest)
		return handlers.NewRegisterHandler(mockAuth).Login(c)
	}

	handleErrorLoginInternalServerError := func(c echo.Context) error {
		mockAuth := func(user requests.LoginAuth) app.AuthService {
			mock := mocks.NewAuthService(t)
			mock.On("Login", user).Return(domain.User{}, "", db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(userMockRequest)
		return handlers.NewRegisterHandler(mockAuth).Login(c)
	}

	handleErrorAccessTokenCreateInternalServer := func(c echo.Context) error {
		mockAuth := func(user requests.LoginAuth) app.AuthService {
			mock := mocks.NewAuthService(t)
			mock.On("Login", user).Return(userMockDomain, "refresh", nil).Times(1).
				On("CreateAccessToken", userMockDomain).Return("", int64(0), errors.New("auth service error create access token"))
			return mock
		}(userMockRequest)
		return handlers.NewRegisterHandler(mockAuth).Login(c)
	}

	handleMock := func(c echo.Context) error {
		mockAuth := func() app.AuthService {
			return mocks.NewAuthService(t)
		}()
		return handlers.NewRegisterHandler(mockAuth).Login(c)
	}

	cases := []test_case.TestCase{
		{
			"LoginUser success",
			requestRegister,
			userMockRequest,
			handleSuccessLogin,
			test_case.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "{\"accessToken\":\"access\",\"refreshToken\":\"refresh\",\"exp\":123}\n"},
		},
		{
			"LoginUser error login no more rows",
			requestRegister,
			userMockRequest,
			handleErrorLoginNoMoreRows,
			test_case.ExpectedResponse{
				StatusCode: 404,
				BodyPart:   "{\"code\":404,\"error\":\"Could not login, user not exists: upper: no more rows in this result set\"}\n"},
		},
		{
			"LoginUser error login no more rows",
			requestRegister,
			userMockRequest,
			handleErrorLoginInternalServerError,
			test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not login user: upper: collection does not exist\"}\n"},
		},
		{
			"LoginUser error login internal server error",
			requestRegister,
			userMockRequest,
			handleErrorAccessTokenCreateInternalServer,
			test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Unauthorized could not create access token: auth service error create access token\"}\n"},
		},
		{
			"Error decode user data",
			requestRegister,
			"",
			handleMock,
			test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not decode user data\"}\n"},
		},
		{
			"Error validate user data",
			requestRegister,
			requests.LoginAuth{
				Email:    "user@email.com",
				Password: "",
			},
			handleMock,
			test_case.ExpectedResponse{
				StatusCode: 422,
				BodyPart:   "{\"code\":422,\"error\":\"Could not validate user data\"}\n"},
		},
	}
	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			c, recorder := test_case.PrepareContextFromTestCase(test)
			c.Set("user", test_case.Token())

			if assert.NoError(t, test.HandlerFunc(c)) {
				assert.Contains(t, recorder.Body.String(), test.Expected.BodyPart)
				assert.Equal(t, recorder.Code, test.Expected.StatusCode)
			}
		})
	}
}
