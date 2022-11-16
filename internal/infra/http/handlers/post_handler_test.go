package handlers_test

import (
	"github.com/golang-jwt/jwt"
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

const postID = "1"

var requestPostMock = requests.PostRequest{
	Title: "title",
	Body:  "body",
}
var returnDomainPostMock = domain.Post{
	UserID: 1,
	ID:     1,
	Title:  "title",
	Body:   "body",
}

func TestWalkPostSuccess(t *testing.T) {
	requestGet := test_case.Request{
		Method: http.MethodGet,
		Url:    "/post/" + postID,
		PathParam: &test_case.PathParam{
			Name:  "id",
			Value: postID,
		},
	}

	requestSave := test_case.Request{
		Method: http.MethodPost,
		Url:    "/save",
	}

	requestUpdate := test_case.Request{
		Method: http.MethodPut,
		Url:    "/update/" + postID,
		PathParam: &test_case.PathParam{
			Name:  "id",
			Value: postID,
		},
	}

	requestDelete := test_case.Request{
		Method: http.MethodDelete,
		Url:    "/delete/" + postID,
		PathParam: &test_case.PathParam{
			Name:  "id",
			Value: postID,
		},
	}

	handleFuncGet := func(c echo.Context) error {
		mock := func(id int64) app.PostService {
			mock := mocks.NewPostService(t)
			mock.On("GetPost", id).Return(returnDomainPostMock, nil).Times(1)
			return mock
		}(1)
		return handlers.NewPostHandler(mock).GetPost(c)
	}

	handleFuncSave := func(c echo.Context) error {
		mock := func(r requests.PostRequest, token *jwt.Token) app.PostService {
			mock := mocks.NewPostService(t)
			mock.On("SavePost", requestPostMock, token).Return(returnDomainPostMock, nil).Times(1)
			return mock
		}(requestPostMock, test_case.Token())
		return handlers.NewPostHandler(mock).SavePost(c)
	}

	handleFuncUpdate := func(c echo.Context) error {
		mock := func(r requests.PostRequest, id int64) app.PostService {
			mock := mocks.NewPostService(t)
			mock.On("UpdatePost", requestPostMock, id).Return(returnDomainPostMock, nil).Times(1)
			return mock
		}(requestPostMock, 1)
		return handlers.NewPostHandler(mock).UpdatePost(c)
	}

	handleFuncDelete := func(c echo.Context) error {
		mock := func(id int64) app.PostService {
			mock := mocks.NewPostService(t)
			mock.On("DeletePost", id).Return(nil).Times(1)
			return mock
		}(1)
		return handlers.NewPostHandler(mock).DeletePost(c)
	}

	cases := []test_case.TestCase{
		{
			"GetPost Success",
			requestGet,
			"",
			handleFuncGet,
			test_case.ExpectedResponse{StatusCode: 200, BodyPart: "{\"id\":1,\"user_id\":1,\"title\":\"title\",\"body\":\"body\"}\n"},
		},
		{
			"SavePost Success",
			requestSave,
			requestPostMock,
			handleFuncSave,
			test_case.ExpectedResponse{
				StatusCode: 201,
				BodyPart:   "{\"id\":1,\"user_id\":1,\"title\":\"title\",\"body\":\"body\"}\n"},
		},
		{
			"UpdatePost Success",
			requestUpdate,
			requestPostMock,
			handleFuncUpdate,
			test_case.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "{\"id\":1,\"user_id\":1,\"title\":\"title\",\"body\":\"body\"}\n"},
		},
		{
			"DeletePost Success",
			requestDelete,
			"",
			handleFuncDelete,
			test_case.ExpectedResponse{
				StatusCode: 200,
			},
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
