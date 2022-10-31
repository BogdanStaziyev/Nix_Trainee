package app

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
	"trainee/config"
	"trainee/internal/domain"
)

const (
	refresh = 144
	access  = 2
)

type AuthService interface {
	Register(user domain.User) (domain.User, error)
	Login(user domain.User) (domain.User, string, error)
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
		log.Println("invalid credentials user exist")
		return domain.User{}, errors.New("invalid credentials user exist")
		//todo wraper
	} else if !errors.Is(err, db.ErrNoMoreRows) {
		//todo wraper
		log.Println(err)
		return domain.User{}, err
	}
	user, err = a.userService.Save(user)
	if err != nil {
		//todo wraper
		log.Println(err)
		return domain.User{}, err
	}
	return user, nil
}

func (a authService) Login(user domain.User) (domain.User, string, error) {
	u, err := a.userService.FindByEmail(user.Email)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			//todo wraper
			log.Println(err)
		}
		log.Println(err)
		return domain.User{}, "", err
	}
	valid := a.checkPasswordHash(user.Password, u.Password)
	if !valid {
		return domain.User{}, "", err
	}
	token, err := a.CreateRefreshToken(u)
	return u, token, err
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
		return "", err
	}
	return token, err
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
		return "", 0, err
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
