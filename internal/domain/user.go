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
