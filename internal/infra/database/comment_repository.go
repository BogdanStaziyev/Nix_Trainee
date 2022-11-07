package database

import (
	"fmt"
	"github.com/upper/db/v4"
	"time"
	"trainee/internal/domain"
	"trainee/internal/infra/http/requests"
)

const CommentTable = "commentses"

type comments struct {
	ID          int64      `db:"id,omitempty"`
	PostID      int64      `db:"post_id"`
	Name        string     `db:"name"`
	Email       string     `db:"email"`
	Body        string     `db:"body"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

//go:generate mockery --dir . --name CommentRepo --output ./mock
type CommentRepo interface {
	SaveComment(comment domain.Comment) (domain.Comment, error)
	GetComment(id int64) (domain.Comment, error)
	UpdateComment(commentRequest requests.CommentRequest, id int64) (domain.Comment, error)
	DeleteComment(id int64) error
	GetCommentsByPostID(postID int64) ([]domain.Comment, error)
}

type commentsRepository struct {
	coll db.Collection
}

func NewCommentRepository(dbSession db.Session) CommentRepo {
	return commentsRepository{
		coll: dbSession.Collection(CommentTable),
	}
}

func (r commentsRepository) SaveComment(comment domain.Comment) (domain.Comment, error) {
	commentsDB := r.mapCommentDBModel(comment)
	commentsDB.CreatedDate = time.Now()
	commentsDB.UpdatedDate = time.Now()
	err := r.coll.InsertReturning(&commentsDB)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("CommentrepositoryCreate: %w", err)
	}
	return r.mapCommentDbModelToDomain(commentsDB), nil
}

func (r commentsRepository) GetComment(id int64) (domain.Comment, error) {
	var comment comments

	err := r.coll.Find(db.Cond{
		"id":           id,
		"deleted_date": nil,
	}).One(&comment)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("CommentRepository GetComment: %w", err)
	}
	return r.mapCommentDbModelToDomain(comment), nil
}

func (r commentsRepository) UpdateComment(commentRequest requests.CommentRequest, id int64) (domain.Comment, error) {
	err := r.coll.Find(db.Cond{
		"id":           id,
		"deleted_date": nil,
	}).Update(map[string]interface{}{
		"updated_date": time.Now(),
		"body":         commentRequest.Body,
	})
	if err != nil {
		return domain.Comment{}, fmt.Errorf("CommentRepository UpdateComment: %w", err)
	}
	return r.GetComment(id)
}

func (r commentsRepository) DeleteComment(id int64) error {
	_, err := r.GetComment(id)
	if err != nil {
		return err
	}
	return r.coll.Find(db.Cond{
		"id":           id,
		"deleted_date": nil,
	}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r commentsRepository) GetCommentsByPostID(postID int64) ([]domain.Comment, error) {
	var comment []comments

	err := r.coll.Find(db.Cond{"post_id": postID}).All(&comment)
	if err != nil {
		return []domain.Comment{}, err
	}
	return r.mapCommentCollection(comment), nil
}

func (r commentsRepository) mapCommentDBModel(comment domain.Comment) comments {
	return comments{
		ID:     comment.ID,
		PostID: comment.PostID,
		Name:   comment.Name,
		Email:  comment.Email,
		Body:   comment.Body,
	}
}

func (r commentsRepository) mapCommentDbModelToDomain(comment comments) domain.Comment {
	return domain.Comment{
		ID:          comment.ID,
		PostID:      comment.PostID,
		Name:        comment.Name,
		Email:       comment.Email,
		Body:        comment.Body,
		CreatedDate: comment.CreatedDate,
		DeletedDate: comment.DeletedDate,
		UpdatedDate: comment.UpdatedDate,
	}
}

func (r commentsRepository) mapCommentCollection(comment []comments) []domain.Comment {
	var result []domain.Comment
	for _, coll := range comment {
		newComment := r.mapCommentDbModelToDomain(coll)
		result = append(result, newComment)
	}
	return result
}
