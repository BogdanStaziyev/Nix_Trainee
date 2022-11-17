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
	commentIDError = "a"
	commentID      = "2"
)

var requestCommentMock = requests.CommentRequest{
	Body: "Test body",
}

var returnDomainCommentMock = domain.Comment{
	ID:     2,
	PostID: 1,
	Name:   "Name",
	Email:  "test@mail.com",
	Body:   "Test body",
}

var requestSaveComment = test_case.Request{
	Method: http.MethodGet,
	Url:    "/save/" + postID,
	PathParam: &test_case.PathParam{
		Name:  "post_id",
		Value: postID,
	},
}

var requestGetComment = test_case.Request{
	Method: http.MethodGet,
	Url:    "/comment/" + commentID,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: commentID,
	},
}

var requestUpdateComment = test_case.Request{
	Method: http.MethodPut,
	Url:    "/update/" + commentID,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: commentID,
	},
}

var requestDeleteComment = test_case.Request{
	Method: http.MethodDelete,
	Url:    "/delete/" + commentID,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: commentID,
	},
}

var requestSaveCommentError = test_case.Request{
	Method: http.MethodGet,
	Url:    "/save/" + postIDError,
	PathParam: &test_case.PathParam{
		Name:  "post_id",
		Value: postIDError,
	},
}

var requestGetCommentError = test_case.Request{
	Method: http.MethodGet,
	Url:    "/comment/" + commentIDError,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: commentIDError,
	},
}

var requestUpdateCommentError = test_case.Request{
	Method: http.MethodPut,
	Url:    "/update/" + commentIDError,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: commentIDError,
	},
}

var requestDeleteCommentError = test_case.Request{
	Method: http.MethodDelete,
	Url:    "/delete/" + commentIDError,
	PathParam: &test_case.PathParam{
		Name:  "id",
		Value: commentIDError,
	},
}

func TestWalkCommentSuccess(t *testing.T) {
	handleFuncGet := func(c echo.Context) error {
		mock := func(id int64) app.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("GetComment", id).Return(returnDomainCommentMock, nil).Times(1)
			return mock
		}(2)
		return handlers.NewCommentHandler(mock).GetComment(c)
	}

	handleFuncSave := func(c echo.Context) error {
		mock := func(r requests.CommentRequest, token *jwt.Token, postID int64) app.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("SaveComment", r, postID, token).Return(returnDomainCommentMock, nil).Times(1)
			return mock
		}(requestCommentMock, test_case.Token(), 1)
		return handlers.NewCommentHandler(mock).SaveComment(c)
	}

	handleFuncUpdate := func(c echo.Context) error {
		mock := func(r requests.CommentRequest, id int64) app.CommentService {
			mock := mocks.NewCommentService(t)
			returnDomainCommentMock.Body = r.Body
			mock.On("UpdateComment", r, id).Return(returnDomainCommentMock, nil).Times(1)
			return mock
		}(requests.CommentRequest{Body: "Update body"}, 2)
		return handlers.NewCommentHandler(mock).UpdateComment(c)
	}

	handleFuncDelete := func(c echo.Context) error {
		mock := func(id int64) app.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("DeleteComment", id).Return(nil).Times(1)
			return mock
		}(2)
		return handlers.NewCommentHandler(mock).DeleteComment(c)
	}

	cases := []test_case.TestCase{
		{
			"GetComment Success",
			requestGetComment,
			"",
			handleFuncGet,
			test_case.ExpectedResponse{StatusCode: 200, BodyPart: "{\"id\":2,\"post_id\":1,\"name\":\"Name\",\"email\":\"test@mail.com\",\"body\":\"Test body\"}\n"},
		},
		{
			"SaveComment Success",
			requestSaveComment,
			requestCommentMock,
			handleFuncSave,
			test_case.ExpectedResponse{
				StatusCode: 201,
				BodyPart:   "{\"id\":2,\"post_id\":1,\"name\":\"Name\",\"email\":\"test@mail.com\",\"body\":\"Test body\"}\n"},
		},
		{
			"UpdateComment Success",
			requestUpdateComment,
			requests.CommentRequest{Body: "Update body"},
			handleFuncUpdate,
			test_case.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "{\"id\":2,\"post_id\":1,\"name\":\"Name\",\"email\":\"test@mail.com\",\"body\":\"Update body\"}\n"},
		},
		{
			"DeleteComment Success",
			requestDeleteComment,
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

func TestWalkCommentsDecodeValidationParsingErrors(t *testing.T) {
	handleFuncGet := func(c echo.Context) error {
		mock := func() app.CommentService {
			mock := mocks.NewCommentService(t)
			return mock
		}()
		return handlers.NewCommentHandler(mock).GetComment(c)
	}

	handleFuncUpdate := func(c echo.Context) error {
		mock := func() app.CommentService {
			mock := mocks.NewCommentService(t)
			return mock
		}()
		return handlers.NewCommentHandler(mock).UpdateComment(c)
	}

	handleFuncDelete := func(c echo.Context) error {
		mock := func() app.CommentService {
			mock := mocks.NewCommentService(t)
			return mock
		}()
		return handlers.NewCommentHandler(mock).DeleteComment(c)
	}

	handleFuncSave := func(c echo.Context) error {
		mock := func() app.CommentService {
			mock := mocks.NewCommentService(t)
			return mock
		}()
		return handlers.NewCommentHandler(mock).SaveComment(c)
	}

	cases := []test_case.TestCase{
		{
			"UpdateComment comment data error",
			requestUpdateComment,
			"",
			handleFuncUpdate,
			test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not decode comment data\"}\n"},
		},
		{
			"SaveComment comment error",
			requestSaveComment,
			"",
			handleFuncSave,
			test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not decode comment data\"}\n"},
		},
		{
			"UpdateComment validate error",
			requestUpdateComment,
			requests.PostRequest{
				Title: "title",
			},
			handleFuncUpdate,
			test_case.ExpectedResponse{
				StatusCode: 422,
				BodyPart:   "{\"code\":422,\"error\":\"Could not validate comment data\"}\n"},
		},
		{
			"SaveComment validate error",
			requestSaveComment,
			requests.PostRequest{
				Title: "title",
			},
			handleFuncSave,
			test_case.ExpectedResponse{
				StatusCode: 422,
				BodyPart:   "{\"code\":422,\"error\":\"Could not validate comment data\"}\n"},
		},
		{
			"SaveComment parse path param Error",
			requestSaveCommentError,
			requestCommentMock,
			handleFuncSave,
			test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not parse post ID\"}\n"},
		},
		{
			"GetComment parse path param Error",
			requestGetCommentError,
			"",
			handleFuncGet,
			test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not parse comment ID\"}\n"},
		},
		{
			"UpdateComment parse path param Error",
			requestUpdateCommentError,
			requestPostMock,
			handleFuncUpdate,
			test_case.ExpectedResponse{
				StatusCode: 400,
				BodyPart:   "{\"code\":400,\"error\":\"Could not parse comment ID\"}\n"},
		},
		{
			"DeleteComment parse path param Error",
			requestDeleteCommentError,
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

func TestWalkCommentServiceErrors(t *testing.T) {
	handleFuncGetNotFound := func(c echo.Context) error {
		mock := func(id int64) app.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("GetComment", id).Return(domain.Comment{}, db.ErrNoMoreRows).Times(1)
			return mock
		}(2)
		return handlers.NewCommentHandler(mock).GetComment(c)
	}

	handleFuncUpdateNotFound := func(c echo.Context) error {
		mock := func(r requests.CommentRequest, id int64) app.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("UpdateComment", r, id).Return(domain.Comment{}, db.ErrNoMoreRows).Times(1)
			return mock
		}(requestCommentMock, 2)
		return handlers.NewCommentHandler(mock).UpdateComment(c)
	}

	handleFuncDeleteNotFound := func(c echo.Context) error {
		mock := func(id int64) app.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("DeleteComment", id).Return(db.ErrNoMoreRows).Times(1)
			return mock
		}(2)
		return handlers.NewCommentHandler(mock).DeleteComment(c)
	}
	handleFuncSaveNotFound := func(c echo.Context) error {
		mock := func(r requests.CommentRequest, p int64, token *jwt.Token) app.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("SaveComment", r, p, token).Return(domain.Comment{}, db.ErrNoMoreRows).Times(1)
			return mock
		}(requestCommentMock, 1, test_case.Token())
		return handlers.NewCommentHandler(mock).SaveComment(c)
	}
	handleFuncGetInternalServerError := func(c echo.Context) error {
		mock := func(id int64) app.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("GetComment", id).Return(domain.Comment{}, db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(2)
		return handlers.NewCommentHandler(mock).GetComment(c)
	}

	handleFuncUpdateInternalServerError := func(c echo.Context) error {
		mock := func(r requests.CommentRequest, id int64) app.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("UpdateComment", r, id).Return(domain.Comment{}, db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(requestCommentMock, 2)
		return handlers.NewCommentHandler(mock).UpdateComment(c)
	}

	handleFuncDeleteInternalServerError := func(c echo.Context) error {
		mock := func(id int64) app.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("DeleteComment", id).Return(db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(2)
		return handlers.NewCommentHandler(mock).DeleteComment(c)
	}

	handleFuncSaveInternalServerError := func(c echo.Context) error {
		mock := func(r requests.CommentRequest, p int64, token *jwt.Token) app.CommentService {
			mock := mocks.NewCommentService(t)
			mock.On("SaveComment", r, p, token).Return(domain.Comment{}, db.ErrCollectionDoesNotExist).Times(1)
			return mock
		}(requestCommentMock, 1, test_case.Token())
		return handlers.NewCommentHandler(mock).SaveComment(c)
	}

	cases := []test_case.TestCase{
		{
			"GetComment NoMoreRows",
			requestGetComment,
			"",
			handleFuncGetNotFound,
			test_case.ExpectedResponse{StatusCode: 404, BodyPart: "{\"code\":404,\"error\":\"Could not get comment: upper: no more rows in this result set\"}\n"},
		},
		{
			"UpdateComment NoMoreRows",
			requestUpdateComment,
			requestCommentMock,
			handleFuncUpdateNotFound,
			test_case.ExpectedResponse{
				StatusCode: 404,
				BodyPart:   "{\"code\":404,\"error\":\"Could not update comment: upper: no more rows in this result set\"}\n"},
		},
		{
			"DeleteComment NoMoreRows",
			requestDeleteComment,
			"",
			handleFuncDeleteNotFound,
			test_case.ExpectedResponse{
				StatusCode: 404,
			},
		},
		{
			"UpdateComment NoMoreRows",
			requestSaveComment,
			requestCommentMock,
			handleFuncSaveNotFound,
			test_case.ExpectedResponse{
				StatusCode: 404,
				BodyPart:   "{\"code\":404,\"error\":\"Could not save new comment: upper: no more rows in this result set\"}\n"},
		},
		{
			"GetComment InternalServerError",
			requestGetComment,
			"",
			handleFuncGetInternalServerError,
			test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not get comment: upper: collection does not exist\"}\n"},
		},
		{
			"UpdateComment InternalServerError",
			requestUpdateComment,
			requestCommentMock,
			handleFuncUpdateInternalServerError,
			test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not update comment: upper: collection does not exist\"}\n"},
		},
		{
			"DeleteComment InternalServerError",
			requestDeleteComment,
			"",
			handleFuncDeleteInternalServerError,
			test_case.ExpectedResponse{
				StatusCode: 500,
			},
		},
		{
			"SaveComment InternalServerError",
			requestSaveComment,
			requestCommentMock,
			handleFuncSaveInternalServerError,
			test_case.ExpectedResponse{
				StatusCode: 500,
				BodyPart:   "{\"code\":500,\"error\":\"Could not save new comment: upper: collection does not exist\"}\n"},
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
