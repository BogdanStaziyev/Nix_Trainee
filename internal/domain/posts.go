package domain

import "time"

type Post struct {
	UserID      int64
	ID          int64
	Title       string
	Body        string
	Comment     []Comment
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}
