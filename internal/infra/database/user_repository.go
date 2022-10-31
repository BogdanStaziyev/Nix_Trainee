package database

import (
	"github.com/upper/db/v4"
	"time"
	"trainee/internal/domain"
)

const UsersTable = "users"

type user struct {
	ID          int64      `db:"id,omitempty"`
	Email       string     `db:"email"`
	Name        string     `db:"name"`
	Password    string     `db:"password,omitempty"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type UserRepo interface {
	Save(user domain.User) (domain.User, error)
	FindByEmail(email string) (domain.User, error)
	FindByID(id int64) (domain.User, error)
	Delete(id int64) error
}

type userRepo struct {
	coll db.Collection
}

func NewUSerRepo(dbSession db.Session) UserRepo {
	return userRepo{
		coll: dbSession.Collection(UsersTable),
	}
}

func (u userRepo) Save(user domain.User) (domain.User, error) {
	domainUser := u.mapDomainToModel(user)
	domainUser.CreatedDate = time.Now()
	domainUser.UpdatedDate = time.Now()
	err := u.coll.InsertReturning(&domainUser)
	if err != nil {
		return domain.User{}, err
	}
	return u.mapModelToDomain(domainUser), nil
}

func (u userRepo) FindByEmail(email string) (domain.User, error) {
	var domainUser user

	err := u.coll.Find(db.Cond{
		"email":       email,
		"delete_date": nil,
	}).One(&domainUser)
	if err != nil {
		return domain.User{}, err
	}
	return u.mapModelToDomain(domainUser), nil
}

func (u userRepo) FindByID(id int64) (domain.User, error) {
	var domainUser user

	err := u.coll.Find(db.Cond{
		"id":           id,
		"deleted_date": nil,
	}).One(&domainUser)
	if err != nil {
		return domain.User{}, err
	}
	return u.mapModelToDomain(domainUser), nil
}

func (u userRepo) Delete(id int64) error {
	return u.coll.Find(db.Cond{
		"id":           id,
		"deleted_date": nil,
	}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (u userRepo) mapDomainToModel(d domain.User) user {
	return user{
		ID:       d.ID,
		Email:    d.Email,
		Password: d.Password,
		Name:     d.Name,
	}
}

func (u userRepo) mapModelToDomain(d user) domain.User {
	return domain.User{
		ID:          d.ID,
		Email:       d.Email,
		Password:    d.Password,
		Name:        d.Name,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}
