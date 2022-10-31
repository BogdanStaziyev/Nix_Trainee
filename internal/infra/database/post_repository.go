package database

import (
	"fmt"
	"github.com/upper/db/v4"
	"time"
	"trainee/internal/domain"
)

const PostTable = "posts"

type posts struct {
	UserID      int64      `db:"user_id"`
	ID          int64      `db:"id,omitempty"`
	Title       string     `db:"title"`
	Body        string     `db:"body"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

//go:generate mockery --dir . --name PostRepo --output ./mock
type PostRepo interface {
	SavePost(post domain.Post) (domain.Post, error)
	GetPost(id int64) (domain.Post, error)
	GetPostsByUser(userID int64) ([]domain.Post, error)
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
	postDB.CreatedDate = time.Now()
	postDB.UpdatedDate = time.Now()
	err := r.coll.InsertReturning(&postDB)
	if err != nil {
		return domain.Post{}, fmt.Errorf("PostRepository Create: %w", err)
	}
	return r.mapPostDbModelToDomain(postDB), nil
}

func (r postsRepository) GetPost(id int64) (domain.Post, error) {
	var post posts

	err := r.coll.Find(db.Cond{
		"id":           id,
		"deleted_date": nil,
	}).One(&post)
	if err != nil {
		return domain.Post{}, fmt.Errorf("PostRepository GetPost: %w", err)
	}
	return r.mapPostDbModelToDomain(post), nil
}

func (r postsRepository) UpdatePost(post domain.Post) (domain.Post, error) {
	updatePost := r.mapPostDBModel(post)
	updatePost.UpdatedDate = time.Now()

	err := r.coll.Find(db.Cond{"id": updatePost.ID}).Update(&updatePost)
	if err != nil {
		return domain.Post{}, fmt.Errorf("PostRepository UpdatePost: %w", err)
	}

	//err = r.coll.Find(db.Cond{"id": updatePost.ID}).One(&updatePost)
	//if err != nil {
	//	return domain.Post{}, err
	//}
	return r.mapPostDbModelToDomain(updatePost), nil
}

func (r postsRepository) DeletePost(id int64) error {
	return r.coll.Find(db.Cond{
		"id":           id,
		"deleted_date": nil,
	}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r postsRepository) GetPostsByUser(userID int64) ([]domain.Post, error) {
	var post []posts

	err := r.coll.Find(db.Cond{"user_id": userID}).All(&post)
	if err != nil {
		return []domain.Post{}, err
	}
	return r.mapPostCollection(post), nil

}

func (r postsRepository) mapPostDBModel(p domain.Post) posts {
	return posts{
		UserID: p.UserID,
		ID:     p.ID,
		Title:  p.Title,
		Body:   p.Body,
	}
}

func (r postsRepository) mapPostDbModelToDomain(p posts) domain.Post {
	return domain.Post{
		UserID:      p.UserID,
		ID:          p.ID,
		Title:       p.Title,
		Body:        p.Body,
		CreatedDate: p.CreatedDate,
		UpdatedDate: p.UpdatedDate,
		DeletedDate: p.DeletedDate,
	}
}

func (r postsRepository) mapPostCollection(post []posts) []domain.Post {
	var result []domain.Post
	for _, coll := range post {
		newPost := r.mapPostDbModelToDomain(coll)
		result = append(result, newPost)
	}
	return result
}
