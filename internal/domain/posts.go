package domain

type Post struct {
	UserId int    `json:"user_id" example:"1" validate:"required"`
	Id     int    `json:"id,omitempty" example:"1"`
	Title  string `json:"title" example:"Lorem ipsum" validate:"required"`
	Body   string `json:"body" example:"Lorem ipsum" validate:"required"`
}
