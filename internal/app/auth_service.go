package app

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
	"trainee/config"
	"trainee/internal/domain"
	"trainee/internal/infra/http/requests"
)

const (
	refresh = 144
	access  = 2
)

//go:generate mockery --dir . --name AuthService --output ./mocks
type AuthService interface {
	Register(user domain.User) (domain.User, error)
	Login(user requests.LoginAuth) (domain.User, string, error)
	CreateRefreshToken(user domain.User) (string, error)
	CreateAccessToken(user domain.User) (string, int64, error)
}

type authService struct {
	userService UserService
	config      config.Configuration
}

func NewAuthService(us UserService, cf config.Configuration) AuthService {
	return authService{
		userService: us,
		config:      cf,
	}
}

func (a authService) Register(user domain.User) (domain.User, error) {
	_, err := a.userService.FindByEmail(user.Email)
	if err == nil {
		return domain.User{}, fmt.Errorf("auth service error register invalid credentials user exist: %w", err)
	} else if !errors.Is(err, db.ErrNoMoreRows) {
		return domain.User{}, fmt.Errorf("auth service error register: %w", err)
	}
	user, err = a.userService.Save(user)
	if err != nil {
		return domain.User{}, fmt.Errorf("auth service error register save user: %w", err)
	}
	return user, nil
}

func (a authService) Login(user requests.LoginAuth) (domain.User, string, error) {
	u, err := a.userService.FindByEmail(user.Email)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return domain.User{}, "", fmt.Errorf("auth service error login, invalid credentials user not exist: %w", err)
		}
		return domain.User{}, "", fmt.Errorf("auth service error login user invalid email or password: %w", err)
	}
	valid := a.checkPasswordHash(user.Password, u.Password)
	if !valid {
		return domain.User{}, "", fmt.Errorf("auth service error login user invalid email or password: %w", err)
	}
	token, err := a.CreateRefreshToken(u)
	if err != nil {
		return domain.User{}, "", fmt.Errorf("auth service error login: %w", err)
	}
	return u, token, nil
}

func (a authService) CreateRefreshToken(user domain.User) (string, error) {
	claimsRefresh := JwtRefreshClaim{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * refresh).Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)

	token, err := refreshToken.SignedString([]byte(a.config.RefreshSecret))
	if err != nil {
		return "", fmt.Errorf("auth service error create refresh token: %w", err)
	}
	return token, nil
}

func (a authService) CreateAccessToken(user domain.User) (string, int64, error) {
	exp := time.Now().Add(time.Hour * access).Unix()
	claimsAccess := JwtAccessClaim{
		Name: user.Name,
		ID:   user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccess)
	t, err := token.SignedString([]byte(a.config.AccessSecret))
	if err != nil {
		return "", 0, fmt.Errorf("auth service error create access token: %w", err)
	}
	return t, exp, err
}

func (a authService) checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

type JwtAccessClaim struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
	jwt.StandardClaims
}

type JwtRefreshClaim struct {
	ID int64 `json:"id"`
	jwt.StandardClaims
}
