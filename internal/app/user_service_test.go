package app

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	smocks "trainee/internal/app/mocks"
	"trainee/internal/domain"
	"trainee/internal/infra/database"
	repoMocks "trainee/internal/infra/database/mock"
)

func Test_userService_FindByEmail(t *testing.T) {
	tests := []struct {
		name            string
		email           string
		repoConstructor func() database.UserRepo
		want            domain.User
		wantErr         bool
	}{
		{
			"success find user by email",
			"user@gmail.com",
			func() database.UserRepo {
				mock := repoMocks.NewUserRepo(t)
				mock.
					On("FindByEmail", "user@gmail.com").
					Return(domain.User{Email: "user@gmail.com"}, nil).Times(1)
				return mock
			},
			domain.User{
				Email: "user@gmail.com",
			},
			false,
		},
		{
			"error find user by email",
			"user@gmail.com",
			func() database.UserRepo {
				mock := repoMocks.NewUserRepo(t)
				mock.
					On("FindByEmail", "user@gmail.com").
					Return(domain.User{}, errors.New("upper: no more rows in this result set")).Times(1)
				return mock
			},
			domain.User{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := userService{
				userRepo: tt.repoConstructor(),
			}
			user, err := us.FindByEmail(tt.email)
			fmt.Println(user)
			if tt.wantErr {
				assert.Error(t, err)
				require.Equal(t, user, domain.User{})
			} else {
				assert.NoError(t, err)
				require.Equal(t, user, tt.want)
			}
		})
	}
}

func Test_userService_FindByID(t *testing.T) {
	tests := []struct {
		name            string
		id              int64
		repoConstructor func() database.UserRepo
		want            domain.User
		wantErr         bool
	}{
		{
			"success find user by id",
			2,
			func() database.UserRepo {
				mock := repoMocks.NewUserRepo(t)
				mock.
					On("FindByID", int64(2)).
					Return(domain.User{ID: 2}, nil).Times(1)
				return mock
			},
			domain.User{
				ID: 2,
			},
			false,
		},
		{
			"error find user by id",
			2,
			func() database.UserRepo {
				mock := repoMocks.NewUserRepo(t)
				mock.
					On("FindByID", int64(2)).
					Return(domain.User{}, errors.New("upper: no more rows in this result set")).Times(1)
				return mock
			},
			domain.User{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := userService{
				userRepo: tt.repoConstructor(),
			}
			user, err := us.FindByID(tt.id)
			fmt.Println(user)
			if tt.wantErr {
				assert.Error(t, err)
				require.Equal(t, user, domain.User{})
			} else {
				assert.NoError(t, err)
				require.Equal(t, user, tt.want)
			}
		})
	}
}

func Test_userService_Save(t *testing.T) {
	tests := []struct {
		name                 string
		user                 domain.User
		repoConstructor      func(user domain.User) database.UserRepo
		generatorConstructor func(password string) Generator
		want                 domain.User
		wantErr              bool
	}{
		{
			"success save user",
			domain.User{
				Email:    "user@mail.com",
				Name:     "user",
				Password: "1234567890",
			},
			func(user domain.User) database.UserRepo {
				user.Password = "0987654321"
				mock := repoMocks.NewUserRepo(t)
				mock.
					On("Save", user).
					Return(domain.User{
						Email:    user.Email,
						Name:     user.Name,
						Password: user.Password,
					}, nil).Times(1)
				return mock
			},
			func(password string) Generator {
				mock := smocks.NewGenerator(t)
				mock.
					On("GeneratePasswordHash", password).
					Return("0987654321", nil).Times(1)
				return mock
			},
			domain.User{
				Email:    "user@mail.com",
				Name:     "user",
				Password: "0987654321",
			},
			false,
		},
		{
			"error save user",
			domain.User{
				Email:    "user@mail.com",
				Name:     "user",
				Password: "1234567890",
			},
			func(user domain.User) database.UserRepo {
				mock := repoMocks.NewUserRepo(t)
				user.Password = "0987654321"
				mock.
					On("Save", user).
					Return(domain.User{}, errors.New("UserService: user already exist")).Times(1)
				return mock
			},
			func(password string) Generator {
				mock := smocks.NewGenerator(t)
				mock.
					On("GeneratePasswordHash", password).
					Return("0987654321", nil).Times(1)
				return mock
			},
			domain.User{},
			true,
		},
		{
			"error save user",
			domain.User{
				Email:    "user@mail.com",
				Name:     "user",
				Password: "1234567890",
			},
			func(user domain.User) database.UserRepo {
				mock := repoMocks.NewUserRepo(t)
				return mock
			},
			func(password string) Generator {
				mock := smocks.NewGenerator(t)
				mock.
					On("GeneratePasswordHash", password).
					Return(password, errors.New("generator error")).Times(1)
				return mock
			},
			domain.User{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := userService{
				passwordGen: tt.generatorConstructor(tt.user.Password),
				userRepo:    tt.repoConstructor(tt.user),
			}
			user, err := NewUserService(us.userRepo, us.passwordGen).Save(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				require.Equal(t, user, domain.User{})
			} else {
				assert.NoError(t, err)
				require.Equal(t, user, tt.want)
			}
		})
	}
}

func Test_userService_Delete(t *testing.T) {
	tests := []struct {
		name            string
		id              int64
		repoConstructor func(id int64) database.UserRepo
		want            error
		wantErr         bool
	}{
		{
			"user delete ok",
			2,
			func(id int64) database.UserRepo {
				mock := repoMocks.NewUserRepo(t)
				mock.
					On("Delete", id).
					Return(nil).Times(1)
				return mock
			},
			nil,
			false,
		},
		{
			"user delete error",
			2,
			func(id int64) database.UserRepo {
				mock := repoMocks.NewUserRepo(t)
				mock.
					On("Delete", id).
					Return(errors.New("could not delete")).Times(1)
				return mock
			},
			errors.New("could not delete"),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := userService{
				userRepo: tt.repoConstructor(tt.id),
			}
			err := u.Delete(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
