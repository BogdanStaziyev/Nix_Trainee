package app

import (
	"log"
	"trainee/internal/domain"
	"trainee/internal/infra/database"
)

//go:generate mockery --dir . --name UserService --output ./mocks
type UserService interface {
	Save(user domain.User) (domain.User, error)
	FindByEmail(email string) (domain.User, error)
	FindByID(id int64) (domain.User, error)
	Delete(id int64) error
}

type userService struct {
	userRepo    database.UserRepo
	passwordGen Generator
}

func NewUserService(ur database.UserRepo, gs Generator) UserService {
	return userService{
		userRepo:    ur,
		passwordGen: gs,
	}
}

func (u userService) Save(user domain.User) (domain.User, error) {
	var err error

	user.Password, err = u.passwordGen.GeneratePasswordHash(user.Password)
	if err != nil {
		log.Printf("UserService: %s", err)
		return domain.User{}, err
	}

	saveUser, err := u.userRepo.Save(user)
	if err != nil {
		log.Printf("UserService: %s", err)
		return domain.User{}, err
	}
	return saveUser, err
}

func (u userService) FindByEmail(email string) (domain.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		log.Println(err)
		return domain.User{}, err
	}
	return user, err
}

func (u userService) FindByID(id int64) (domain.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		log.Println(err)
		return domain.User{}, err
	}
	return user, err
}

func (u userService) Delete(id int64) error {
	err := u.userRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
