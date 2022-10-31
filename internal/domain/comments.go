package domain

import "time"

type Comment struct {
	ID          int64
	PostID      int64
	Name        string
	Email       string
	Body        string
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

//ID          int64  `json:"id" example:"1"`
//PostID      int64  `json:"post_id" example:"3" validate:"required"`
//Name        string `json:"name" example:"Bohdan" validate:"required"`
//Email       string `json:"email" example:"example@mail.com" validate:"required"`
//Body        string `json:"body" example:"lorem ipsum" validate:"required"`
