package app

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"trainee/internal/domain"
	"trainee/internal/infra/database"
	mocks "trainee/internal/infra/database/mock"
	"trainee/internal/infra/http/requests"
)

func token() *jwt.Token {
	exp := time.Now().Add(time.Hour * access).Unix()
	claimsAccess := &JwtTokenClaim{
		Name: "Name",
		ID:   int64(1),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	tokenReturn := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccess)
	return tokenReturn
}

func Test_postService_SavePost(t *testing.T) {
	tests := []struct {
		name            string
		post            domain.Post
		postRequest     requests.PostRequest
		repoConstructor func(post domain.Post) database.PostRepo
		token           *jwt.Token
		want            domain.Post
		wantErr         bool
	}{
		{
			"Success create post",
			domain.Post{
				UserID: 1,
				Title:  "Title",
				Body:   "Body",
			},
			requests.PostRequest{
				Title: "Title",
				Body:  "Body",
			},
			func(post domain.Post) database.PostRepo {
				mock := mocks.NewPostRepo(t)
				mock.
					On("SavePost", post).
					Return(post, nil)
				return mock
			},
			token(),
			domain.Post{
				UserID: 1,
				Title:  "Title",
				Body:   "Body",
			},
			false,
		},
		{
			"Error create post",
			domain.Post{
				UserID: 1,
				Title:  "Title",
				Body:   "Body",
			},
			requests.PostRequest{
				Title: "Title",
				Body:  "Body",
			},
			func(post domain.Post) database.PostRepo {
				mock := mocks.NewPostRepo(t)
				mock.
					On("SavePost", post).
					Return(domain.Post{}, errors.New("error")).Times(1)
				return mock
			},
			token(),
			domain.Post{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := postService{
				repo: tt.repoConstructor(tt.post),
			}
			post, err := NewPostService(s.repo).SavePost(tt.postRequest, tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				require.Equal(t, post, tt.want)
			} else {
				assert.NoError(t, err)
				require.Equal(t, post, tt.want)
			}
		})
	}
}

func Test_postService_GetPost(t *testing.T) {
	tests := []struct {
		name            string
		id              int64
		repoConstructor func(id int64) database.PostRepo
		want            domain.Post
		wantErr         bool
	}{
		{
			"Success get post",
			2,
			func(id int64) database.PostRepo {
				mock := mocks.NewPostRepo(t)
				mock.
					On("GetPost", id).
					Return(domain.Post{ID: id}, nil)
				return mock
			},
			domain.Post{
				ID: 2,
			},
			false,
		},
		{
			"Error get post",
			2,
			func(id int64) database.PostRepo {
				mock := mocks.NewPostRepo(t)
				mock.
					On("GetPost", id).
					Return(domain.Post{}, errors.New("error")).Times(1)
				return mock
			},
			domain.Post{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := postService{
				repo: tt.repoConstructor(tt.id),
			}
			post, err := NewPostService(s.repo).GetPost(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				require.Equal(t, post, tt.want)
			} else {
				assert.NoError(t, err)
				require.Equal(t, post, tt.want)
			}
		})
	}
}

func Test_postService_UpdatePost(t *testing.T) {
	tests := []struct {
		name            string
		postRequest     requests.PostRequest
		postID          int64
		repoConstructor func(post domain.Post, postID int64) database.PostRepo
		want            domain.Post
		wantErr         bool
	}{
		{
			"Success update post",
			requests.PostRequest{
				Title: "Title",
				Body:  "Body",
			},
			2,
			func(post domain.Post, id int64) database.PostRepo {
				mock := mocks.NewPostRepo(t)
				mock.
					On("GetPost", id).
					Return(domain.Post{ID: id}, nil).
					On("UpdatePost", post).
					Return(domain.Post{
						ID:    2,
						Title: post.Title,
						Body:  post.Body,
					}, nil)
				return mock
			},
			domain.Post{
				ID:    2,
				Title: "Title",
				Body:  "Body",
			},
			false,
		},
		{
			"Error update post GetPost repo",
			requests.PostRequest{
				Title: "Title",
				Body:  "Body",
			},
			2,
			func(post domain.Post, id int64) database.PostRepo {
				mock := mocks.NewPostRepo(t)
				mock.
					On("GetPost", id).
					Return(domain.Post{}, errors.New("upper: no more rows in this result set"))
				return mock
			},
			domain.Post{},
			true,
		},
		{
			"Error update post UpdatePost repo",
			requests.PostRequest{
				Title: "Title",
				Body:  "Body",
			},
			2,
			func(post domain.Post, id int64) database.PostRepo {
				mock := mocks.NewPostRepo(t)
				mock.
					On("GetPost", id).
					Return(domain.Post{ID: id}, nil).
					On("UpdatePost", post).
					Return(domain.Post{}, errors.New("post repository update post"))
				return mock
			},
			domain.Post{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := domain.Post{
				ID:    tt.postID,
				Title: tt.postRequest.Title,
				Body:  tt.postRequest.Body,
			}
			s := postService{
				repo: tt.repoConstructor(post, tt.postID),
			}
			post, err := NewPostService(s.repo).UpdatePost(tt.postRequest, tt.postID)
			if tt.wantErr {
				assert.Error(t, err)
				require.Equal(t, post, tt.want)
			} else {
				assert.NoError(t, err)
				require.Equal(t, post, tt.want)
			}
		})
	}
}

func Test_postService_DeletePost(t *testing.T) {
	tests := []struct {
		name            string
		postID          int64
		repoConstructor func(postID int64) database.PostRepo
		wantErr         bool
	}{
		{
			"Success delete post",
			2,
			func(id int64) database.PostRepo {
				mock := mocks.NewPostRepo(t)
				mock.
					On("GetPost", id).
					Return(domain.Post{}, nil).
					On("DeletePost", id).
					Return(nil)
				return mock
			},
			false,
		},
		{
			"Error delete post GetPost repo",
			2,
			func(id int64) database.PostRepo {
				mock := mocks.NewPostRepo(t)
				mock.
					On("GetPost", id).
					Return(domain.Post{}, errors.New("upper: no more rows in this result set"))
				return mock
			},
			true,
		},
		{
			"Error delete post DeletePost repo",
			2,
			func(id int64) database.PostRepo {
				mock := mocks.NewPostRepo(t)
				mock.
					On("GetPost", id).
					Return(domain.Post{}, nil).
					On("DeletePost", id).
					Return(errors.New("post repository delete post"))
				return mock
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := postService{
				repo: tt.repoConstructor(tt.postID),
			}
			err := NewPostService(s.repo).DeletePost(tt.postID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_postService_GetPostsByUser(t *testing.T) {
	tests := []struct {
		name            string
		userID          int64
		repoConstructor func(userID int64) database.PostRepo
		want            []domain.Post
		wantErr         bool
	}{
		{
			"Success get posts by user id",
			2,
			func(userID int64) database.PostRepo {
				mock := mocks.NewPostRepo(t)
				mock.
					On("GetPostsByUser", userID).
					Return([]domain.Post{
						{UserID: 2},
						{UserID: 2},
						{UserID: 2},
					}, nil)
				return mock
			},
			[]domain.Post{
				{UserID: 2},
				{UserID: 2},
				{UserID: 2},
			},
			false,
		},
		{
			"Error delete post GetPost repo",
			2,
			func(userID int64) database.PostRepo {
				mock := mocks.NewPostRepo(t)
				mock.
					On("GetPostsByUser", userID).
					Return([]domain.Post{}, errors.New("upper: no more rows in this result set"))
				return mock
			},
			[]domain.Post{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := postService{
				repo: tt.repoConstructor(tt.userID),
			}
			posts, err := NewPostService(s.repo).GetPostsByUser(tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				require.Equal(t, posts, tt.want)
			} else {
				assert.NoError(t, err)
				require.Equal(t, posts, tt.want)
			}
		})
	}
}
