package resources

import "trainee/internal/domain"

type PostDTO struct {
	UserId int    `json:"user_id"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func MapDomainToPostDTO(p domain.Post) PostDTO {
	return PostDTO{
		UserId: p.UserId,
		Id:     p.Id,
		Title:  p.Title,
		Body:   p.Body,
	}
}
