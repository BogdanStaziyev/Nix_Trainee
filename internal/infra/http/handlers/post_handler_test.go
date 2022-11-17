package handlers_test

import (
	"github.com/golang-jwt/jwt"
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

const (
	postID      = "1"
	postIDError = "a"
)

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

var requestGet = test_case.Request{
	Method: http.MethodGet,
	Url:    "/post/" + postID,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: postID,
	},
}

var requestSave = test_case.Request{
	Method: http.MethodPost,
	Url:    "/save",
}

var requestUpdate = test_case.Request{
	Method: http.MethodPut,
	Url:    "/update/" + postID,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: postID,
	},
}

var requestDelete = test_case.Request{
	Method: http.MethodDelete,
	Url:    "/delete/" + postID,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: postID,
	},
}

var requestGetError = test_case.Request{
	Method: http.MethodGet,
	Url:    "/post/" + postIDError,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: postIDError,
	},
}

var requestUpdateError = test_case.Request{
	Method: http.MethodPut,
	Url:    "/update/" + postIDError,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: postIDError,
	},
}

var requestDeleteError = test_case.Request{
	Method: http.MethodDelete,
	Url:    "/delete/" + postIDError,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: postIDError,
	},
}

func TestWalkPostSuccess(t *testing.T) {
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
			test_case.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "{\"id\":1,\"user_id\":1,\"title\":\"title\",\"body\":\"body\"}\n"},
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

func TestWalkPostsDecodeValidationParsingErrors(t *testing.T) {
	handleFuncGet := func(c echo.Context) error {
		mock := func() app.PostService {
			mock := mocks.NewPostService(t)
			return mock
		}()
		return handlers.NewPostHandler(mock).GetPost(c)
	}

	handleFuncUpdate := func(c echo.Context) error {
		mock := func() app.PostService {
			mock := mocks.NewPostService(t)
			return mock
		}()
		return handlers.NewPostHandler(mock).UpdatePost(c)
	}

	handleFuncDelete := func(c echo.Context) error {
		mock := func() app.PostService {
			mock := mocks.NewPostService(t)
			return mock
		}()
		return handlers.NewPostHandler(mock).DeletePost(c)
	}

	handleFuncSave := func(c echo.Context) error {
		mock := func() app.PostService {
			mock := mocks.NewPostService(t)
			return mock
		}()
		return handlers.NewPostHandler(mock).SavePost(c)
	}

	cases := []test_case.TestCase{
		{
			"UpdatePost post data error",
			requestUpdate,
			"",
			handleFuncUpdate,
			test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not decode post data\"}\n"},
		},
		{
			"SavePost post data error",
			requestSave,
			"",
			handleFuncSave,
			test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not decode post data\"}\n"},
		},
		{
			"UpdatePost validate error",
			requestUpdate,
			requests.PostRequest{
				Title: "title",
			},
			handleFuncUpdate,
			test_case.ExpectedResponse{
				StatusCode: 422,
				BodyPart:   "{\"code\":422,\"error\":\"Could not validate post data\"}\n"},
		},
		{
			"SavePost validate error",
			requestSave,
			requests.PostRequest{
				Title: "title",
			},
			handleFuncSave,
			test_case.ExpectedResponse{
				StatusCode: 422,
				BodyPart:   "{\"code\":422,\"error\":\"Could not validate post data\"}\n"},
		},
		{
			"GetPost parse path param Error",
			requestGetError,
			"",
			handleFuncGet,
			test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not parse post ID\"}\n"},
		},
		{
			"UpdatePost parse path param Error",
			requestUpdateError,
			requestPostMock,
			handleFuncUpdate,
			test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not parse post ID\"}\n"},
		},
		{
			"DeletePost parse path param Error",
			requestDeleteError,
			"",
			handleFuncDelete,
			test_case.ExpectedResponse{
				StatusCode: 400,
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

func TestWalkPostServiceErrors(t *testing.T) {
	handleFuncGetNotFound := func(c echo.Context) error {
		mock := func(id int64) app.PostService {
			mock := mocks.NewPostService(t)
			mock.On("GetPost", id).Return(domain.Post{}, db.ErrNoMoreRows).Times(1)
			return mock
		}(1)
		return handlers.NewPostHandler(mock).GetPost(c)
	}

	handleFuncUpdateNotFound := func(c echo.Context) error {
		mock := func(r requests.PostRequest, id int64) app.PostService {
			mock := mocks.NewPostService(t)
			mock.On("UpdatePost", requestPostMock, id).Return(domain.Post{}, db.ErrNoMoreRows).Times(1)
			return mock
		}(requestPostMock, 1)
		return handlers.NewPostHandler(mock).UpdatePost(c)
	}

	handleFuncDeleteNotFound := func(c echo.Context) error {
		mock := func(id int64) app.PostService {
			mock := mocks.NewPostService(t)
			mock.On("DeletePost", id).Return(db.ErrNoMoreRows).Times(1)
			return mock
		}(1)
		return handlers.NewPostHandler(mock).DeletePost(c)
	}
	handleFuncGetInternalServerError := func(c echo.Context) error {
		mock := func(id int64) app.PostService {
			mock := mocks.NewPostService(t)
			mock.On("GetPost", id).Return(domain.Post{}, db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(1)
		return handlers.NewPostHandler(mock).GetPost(c)
	}

	handleFuncUpdateInternalServerError := func(c echo.Context) error {
		mock := func(r requests.PostRequest, id int64) app.PostService {
			mock := mocks.NewPostService(t)
			mock.On("UpdatePost", requestPostMock, id).Return(domain.Post{}, db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(requestPostMock, 1)
		return handlers.NewPostHandler(mock).UpdatePost(c)
	}

	handleFuncDeleteInternalServerError := func(c echo.Context) error {
		mock := func(id int64) app.PostService {
			mock := mocks.NewPostService(t)
			mock.On("DeletePost", id).Return(db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(1)
		return handlers.NewPostHandler(mock).DeletePost(c)
	}

	handleFuncSaveInternalServerError := func(c echo.Context) error {
		mock := func(r requests.PostRequest, token *jwt.Token) app.PostService {
			mock := mocks.NewPostService(t)
			mock.On("SavePost", requestPostMock, token).Return(domain.Post{}, db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(requestPostMock, test_case.Token())
		return handlers.NewPostHandler(mock).SavePost(c)
	}

	cases := []test_case.TestCase{
		{
			"GetPost NoMoreRows",
			requestGet,
			"",
			handleFuncGetNotFound,
			test_case.ExpectedResponse{StatusCode: 404, BodyPart: "{\"code\":404,\"error\":\"Could not get post: upper: no more rows in this result set\"}\n"},
		},
		{
			"UpdatePost NoMoreRows",
			requestUpdate,
			requestPostMock,
			handleFuncUpdateNotFound,
			test_case.ExpectedResponse{
				StatusCode: 404,
				BodyPart:   "{\"code\":404,\"error\":\"Could not get post: upper: no more rows in this result set\"}\n"},
		},
		{
			"DeletePost NoMoreRows",
			requestDelete,
			"",
			handleFuncDeleteNotFound,
			test_case.ExpectedResponse{
				StatusCode: 404,
			},
		},
		{
			"GetPost InternalServerError",
			requestGet,
			"",
			handleFuncGetInternalServerError,
			test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not get post: upper: collection does not exist\"}\n"},
		},
		{
			"UpdatePost InternalServerError",
			requestUpdate,
			requestPostMock,
			handleFuncUpdateInternalServerError,
			test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not get post: upper: collection does not exist\"}\n"},
		},
		{
			"DeletePost InternalServerError",
			requestDelete,
			"",
			handleFuncDeleteInternalServerError,
			test_case.ExpectedResponse{
				StatusCode: 500,
			},
		},
		{
			"SavePost InternalServerError",
			requestSave,
			requestPostMock,
			handleFuncSaveInternalServerError,
			test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not save new post: upper: collection does not exist\"}\n"},
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
