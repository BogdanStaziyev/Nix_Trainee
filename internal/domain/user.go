package domain

import "time"

type User struct {
	ID          int64
	Email       string
	Name        string
	Password    string
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

//type AuthUser struct {
//	Email    string `json:"email" validate:"required"`
//	Password string `json:"password" validate:"required"`
//}
//
//type RegisterUser struct {
//	AuthUser
//	Name string `json:"name" validate:"required"`
//}
//
//type Refresh struct {
//	Token string `json:"token" validate:"required"`
//}

//ID       int64  `json:"id"`
//Email    string `json:"email" validate:"required,email"`
//Name     string `json:"name" validate:"required,gte=3"`
//Password string `json:"password" validate:"required,gte=8"`
//Post     []Post
