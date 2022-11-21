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
			TestName:    "GetPost Success",
			Request:     requestGet,
			RequestBody: "",
			HandlerFunc: handleFuncGet,
			Expected: test_case.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "{\"id\":1,\"user_id\":1,\"title\":\"title\",\"body\":\"body\"}\n"},
		},
		{
			TestName:    "SavePost Success",
			Request:     requestSave,
			RequestBody: requestPostMock,
			HandlerFunc: handleFuncSave,
			Expected: test_case.ExpectedResponse{
				StatusCode: 201,
				BodyPart:   "{\"id\":1,\"user_id\":1,\"title\":\"title\",\"body\":\"body\"}\n"},
		},
		{
			TestName:    "UpdatePost Success",
			Request:     requestUpdate,
			RequestBody: requestPostMock,
			HandlerFunc: handleFuncUpdate,
			Expected: test_case.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "{\"id\":1,\"user_id\":1,\"title\":\"title\",\"body\":\"body\"}\n"},
		},
		{
			TestName:    "DeletePost Success",
			Request:     requestDelete,
			RequestBody: "",
			HandlerFunc: handleFuncDelete,
			Expected: test_case.ExpectedResponse{
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
			TestName:    "UpdatePost post data error",
			Request:     requestUpdate,
			RequestBody: "",
			HandlerFunc: handleFuncUpdate,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not decode post data\"}\n"},
		},
		{
			TestName:    "SavePost post data error",
			Request:     requestSave,
			RequestBody: "",
			HandlerFunc: handleFuncSave,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not decode post data\"}\n"},
		},
		{
			TestName: "UpdatePost validate error",
			Request:  requestUpdate,
			RequestBody: requests.PostRequest{
				Title: "title",
			},
			HandlerFunc: handleFuncUpdate,
			Expected: test_case.ExpectedResponse{
				StatusCode: 422,
				BodyPart:   "{\"code\":422,\"error\":\"Could not validate post data\"}\n"},
		},
		{
			TestName: "SavePost validate error",
			Request:  requestSave,
			RequestBody: requests.PostRequest{
				Title: "title",
			},
			HandlerFunc: handleFuncSave,
			Expected: test_case.ExpectedResponse{
				StatusCode: 422,
				BodyPart:   "{\"code\":422,\"error\":\"Could not validate post data\"}\n"},
		},
		{
			TestName:    "GetPost parse path param Error",
			Request:     requestGetError,
			RequestBody: "",
			HandlerFunc: handleFuncGet,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not parse post ID\"}\n"},
		},
		{
			TestName:    "UpdatePost parse path param Error",
			Request:     requestUpdateError,
			RequestBody: requestPostMock,
			HandlerFunc: handleFuncUpdate,
			Expected: test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not parse post ID\"}\n"},
		},
		{
			TestName:    "DeletePost parse path param Error",
			Request:     requestDeleteError,
			RequestBody: "",
			HandlerFunc: handleFuncDelete,
			Expected: test_case.ExpectedResponse{
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
			TestName:    "GetPost NoMoreRows",
			Request:     requestGet,
			RequestBody: "",
			HandlerFunc: handleFuncGetNotFound,
			Expected:    test_case.ExpectedResponse{StatusCode: 404, BodyPart: "{\"code\":404,\"error\":\"Could not get post: upper: no more rows in this result set\"}\n"},
		},
		{
			TestName:    "UpdatePost NoMoreRows",
			Request:     requestUpdate,
			RequestBody: requestPostMock,
			HandlerFunc: handleFuncUpdateNotFound,
			Expected: test_case.ExpectedResponse{
				StatusCode: 404,
				BodyPart:   "{\"code\":404,\"error\":\"Could not get post: upper: no more rows in this result set\"}\n"},
		},
		{
			TestName:    "DeletePost NoMoreRows",
			Request:     requestDelete,
			RequestBody: "",
			HandlerFunc: handleFuncDeleteNotFound,
			Expected: test_case.ExpectedResponse{
				StatusCode: 404,
			},
		},
		{
			TestName:    "GetPost InternalServerError",
			Request:     requestGet,
			RequestBody: "",
			HandlerFunc: handleFuncGetInternalServerError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not get post: upper: collection does not exist\"}\n"},
		},
		{
			TestName:    "UpdatePost InternalServerError",
			Request:     requestUpdate,
			RequestBody: requestPostMock,
			HandlerFunc: handleFuncUpdateInternalServerError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not get post: upper: collection does not exist\"}\n"},
		},
		{
			TestName:    "DeletePost InternalServerError",
			Request:     requestDelete,
			RequestBody: "",
			HandlerFunc: handleFuncDeleteInternalServerError,
			Expected: test_case.ExpectedResponse{
				StatusCode: 500,
			},
		},
		{
			TestName:    "SavePost InternalServerError",
			Request:     requestSave,
			RequestBody: requestPostMock,
			HandlerFunc: handleFuncSaveInternalServerError,
			Expected: test_case.ExpectedResponse{
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
