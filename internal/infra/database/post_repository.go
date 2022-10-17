package database

import (
	"fmt"
	"github.com/upper/db/v4"
	"trainee/internal/domain"
)

const PostTable = "posts"

type posts struct {
	UserId int    `db:"user_id"`
	Id     int    `db:"id,omitempty"`
	Title  string `db:"title"`
	Body   string `db:"body"`
}

type PostRepo interface {
	SavePost(post domain.Post) (domain.Post, error)
	GetPost(id int64) (domain.Post, error)
	UpdatePost(post domain.Post) (domain.Post, error)
	DeletePost(id int64) error
}

type postsRepository struct {
	coll db.Collection
}

func NewPostRepository(dbSession db.Session) PostRepo {
	return postsRepository{
		coll: dbSession.Collection(PostTable),
	}
}

func (r postsRepository) SavePost(post domain.Post) (domain.Post, error) {
	postDB := r.mapPostDBModel(post)
	err := r.coll.InsertReturning(&postDB)
	if err != nil {
		return domain.Post{}, fmt.Errorf("PostRepository Create: %w", err)
	}
	return r.mapPostDbModelToDomain(postDB), nil
}

func (r postsRepository) GetPost(id int64) (domain.Post, error) {
	var post posts

	err := r.coll.Find("id", id).One(&post)
	if err != nil {
		return domain.Post{}, fmt.Errorf("PostRepository GetPost: %w", err)
	}
	return r.mapPostDbModelToDomain(post), nil
}

func (r postsRepository) UpdatePost(post domain.Post) (domain.Post, error) {
	updatePost := r.mapPostDBModel(post)

	err := r.coll.Find(db.Cond{"id": updatePost.Id}).Update(&updatePost)
	if err != nil {
		return domain.Post{}, fmt.Errorf("PostRepository UpdatePost: %w", err)
	}

	err = r.coll.Find(db.Cond{"id": updatePost.Id}).One(&updatePost)
	if err != nil {
		return domain.Post{}, err
	}
	return r.mapPostDbModelToDomain(updatePost), nil
}

func (r postsRepository) DeletePost(id int64) error {
	err := r.coll.Find("id", id).Delete()
	if err != nil {
		return fmt.Errorf("PostRepository Delete: %w", err)
	}
	return nil
}

func (r postsRepository) mapPostDBModel(p domain.Post) posts {
	return posts{
		UserId: p.UserId,
		Id:     p.Id,
		Title:  p.Title,
		Body:   p.Body,
	}
}

func (r postsRepository) mapPostDbModelToDomain(p posts) domain.Post {
	return domain.Post{
		UserId: p.UserId,
		Id:     p.Id,
		Title:  p.Title,
		Body:   p.Body,
	}
}
