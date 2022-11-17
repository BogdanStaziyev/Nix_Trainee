package handlers_test

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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
	handleFunc := func(c echo.Context) error {
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
		mockUser := func() app.UserService {
			return mocks.NewUserService(t)
		}()
		return handlers.NewRegisterHandler(mockUser, mockAuth).Register(c)
	}
	cases := []test_case.TestCase{
		{
			"RegisterUser success",
			requestRegister,
			userMockRequest,
			handleFunc,
			test_case.ExpectedResponse{
				StatusCode: 201,
				BodyPart:   "{\"id\":1,\"email\":\"user@mail.com\",\"name\":\"Name\"}\n"},
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
