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

const commentID = "2"

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

func TestWalkCommentSuccess(t *testing.T) {
	requestSave := test_case.Request{
		Method: http.MethodGet,
		Url:    "/save/" + postID,
		PathParam: &test_case.PathParam{
			Name:  "post_id",
			Value: postID,
		},
	}

	requestGet := test_case.Request{
		Method: http.MethodGet,
		Url:    "/comment/" + commentID,
		PathParam: &test_case.PathParam{
			Name:  "id",
			Value: commentID,
		},
	}

	requestUpdate := test_case.Request{
		Method: http.MethodPut,
		Url:    "/update/" + commentID,
		PathParam: &test_case.PathParam{
			Name:  "id",
			Value: commentID,
		},
	}

	requestDelete := test_case.Request{
		Method: http.MethodDelete,
		Url:    "/delete/" + commentID,
		PathParam: &test_case.PathParam{
			Name:  "id",
			Value: commentID,
		},
	}

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
			requestGet,
			"",
			handleFuncGet,
			test_case.ExpectedResponse{StatusCode: 200, BodyPart: "{\"id\":2,\"post_id\":1,\"name\":\"Name\",\"email\":\"test@mail.com\",\"body\":\"Test body\"}\n"},
		},
		{
			"SaveComment Success",
			requestSave,
			requestCommentMock,
			handleFuncSave,
			test_case.ExpectedResponse{
				StatusCode: 201,
				BodyPart:   "{\"id\":2,\"post_id\":1,\"name\":\"Name\",\"email\":\"test@mail.com\",\"body\":\"Test body\"}\n"},
		},
		{
			"UpdateComment Success",
			requestUpdate,
			requests.CommentRequest{Body: "Update body"},
			handleFuncUpdate,
			test_case.ExpectedResponse{
				StatusCode: 200,
				BodyPart:   "{\"id\":2,\"post_id\":1,\"name\":\"Name\",\"email\":\"test@mail.com\",\"body\":\"Update body\"}\n"},
		},
		{
			"DeleteComment Success",
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
