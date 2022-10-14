package domain

type Comment struct {
	Id     int64
	PostId int64
	Name   string
	Email  string
	Body   string
}

//todo delete json
