package app

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/upper/db/v4"
	"testing"
	smocks "trainee/internal/app/mocks"
	"trainee/internal/domain"
	"trainee/internal/infra/database"
	rmocks "trainee/internal/infra/database/mock"
	"trainee/internal/infra/http/requests"
)

func Test_commentService_GetComment(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		repo    func(id int64) database.CommentRepo
		want    domain.Comment
		wantErr bool
	}{
		{
			"success get comment",
			2,
			func(id int64) database.CommentRepo {
				mock := rmocks.NewCommentRepo(t)
				mock.
					On("GetComment", id).
					Return(domain.Comment{ID: 2}, nil)
				return mock
			},
			domain.Comment{ID: 2},
			false,
		},
		{
			"error get comment",
			2,
			func(id int64) database.CommentRepo {
				mock := rmocks.NewCommentRepo(t)
				mock.
					On("GetComment", id).
					Return(domain.Comment{}, db.ErrNoMoreRows)
				return mock
			},
			domain.Comment{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := commentService{
				repo: tt.repo(tt.id),
			}
			comment, err := NewCommentService(s.repo, s.us, s.ps).GetComment(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, comment, tt.want)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, comment, tt.want)
			}
		})
	}
}

func Test_commentService_SaveComment(t *testing.T) {
	tests := []struct {
		name           string
		commentRequest requests.CommentRequest
		postID         int64
		token          *jwt.Token
		ps             func(postID int64) PostService
		us             func(id int64) UserService
		repo           func(commentRequest requests.CommentRequest) database.CommentRepo
		want           domain.Comment
		wantErr        bool
	}{
		{
			"success create comment",
			requests.CommentRequest{Body: "body"},
			2,
			token(),
			func(postID int64) PostService {
				mock := smocks.NewPostService(t)
				mock.
					On("GetPost", postID).
					Return(domain.Post{}, nil).Times(1)
				return mock
			},
			func(id int64) UserService {
				mock := smocks.NewUserService(t)
				mock.
					On("FindByID", id).
					Return(domain.User{
						ID:    id,
						Email: "comment@mail.com",
						Name:  "Name",
					}, nil).Times(1)
				return mock
			},
			func(commentRequest requests.CommentRequest) database.CommentRepo {
				mock := rmocks.NewCommentRepo(t)
				domainComment := domain.Comment{
					PostID: 2,
					Name:   "Name",
					Email:  "comment@mail.com",
					Body:   "body",
				}
				mock.
					On("SaveComment", domainComment).
					Return(domainComment, nil)
				return mock
			},
			domain.Comment{
				PostID: 2,
				Name:   "Name",
				Email:  "comment@mail.com",
				Body:   "body",
			},
			false,
		},
		{
			"error post not exist create comment",
			requests.CommentRequest{Body: "body"},
			2,
			token(),
			func(postID int64) PostService {
				mock := smocks.NewPostService(t)
				mock.
					On("GetPost", postID).
					Return(domain.Post{}, db.ErrNoMoreRows).Times(1)
				return mock
			},
			func(id int64) UserService {
				mock := smocks.NewUserService(t)
				return mock
			},
			func(commentRequest requests.CommentRequest) database.CommentRepo {
				mock := rmocks.NewCommentRepo(t)
				return mock
			},
			domain.Comment{},
			true,
		},
		{
			"error user not exist create comment",
			requests.CommentRequest{Body: "body"},
			2,
			token(),
			func(postID int64) PostService {
				mock := smocks.NewPostService(t)
				mock.
					On("GetPost", postID).
					Return(domain.Post{}, nil).Times(1)
				return mock
			},
			func(id int64) UserService {
				mock := smocks.NewUserService(t)
				mock.
					On("FindByID", id).
					Return(domain.User{}, db.ErrNoMoreRows).Times(1)
				return mock
			},
			func(commentRequest requests.CommentRequest) database.CommentRepo {
				mock := rmocks.NewCommentRepo(t)
				return mock
			},
			domain.Comment{},
			true,
		},
		{
			"error database create comment",
			requests.CommentRequest{Body: "body"},
			2,
			token(),
			func(postID int64) PostService {
				mock := smocks.NewPostService(t)
				mock.
					On("GetPost", postID).
					Return(domain.Post{}, nil).Times(1)
				return mock
			},
			func(id int64) UserService {
				mock := smocks.NewUserService(t)
				mock.
					On("FindByID", id).
					Return(domain.User{
						ID:    id,
						Email: "comment@mail.com",
						Name:  "Name",
					}, nil).Times(1)
				return mock
			},
			func(commentRequest requests.CommentRequest) database.CommentRepo {
				mock := rmocks.NewCommentRepo(t)
				domainComment := domain.Comment{
					PostID: 2,
					Name:   "Name",
					Email:  "comment@mail.com",
					Body:   "body",
				}
				mock.
					On("SaveComment", domainComment).
					Return(domain.Comment{}, db.ErrMissingPrimaryKeys)
				return mock
			},
			domain.Comment{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := commentService{
				repo: tt.repo(tt.commentRequest),
				us:   tt.us(tt.token.Claims.(*JwtAccessClaim).ID),
				ps:   tt.ps(tt.postID),
			}
			comment, err := NewCommentService(s.repo, s.us, s.ps).SaveComment(tt.commentRequest, tt.postID, tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, comment, tt.want)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, comment, tt.want)
			}
		})
	}
}
