package database

import (
	"fmt"
	"github.com/upper/db/v4"
	"trainee/internal/domain"
)

const CommentTable = "commentses"

type comments struct {
	ID     int64  `db:"id,omitempty"`
	PostID int64  `db:"post_id"`
	Name   string `db:"name"`
	Email  string `db:"email"`
	Body   string `db:"body"`
}

//go:generate mockery --dir . --name CommentRepo --output ./mock
type CommentRepo interface {
	SaveComment(comment domain.Comment) (domain.Comment, error)
	GetComment(id int64) (domain.Comment, error)
	UpdateComment(comment domain.Comment) (domain.Comment, error)
	DeleteComment(id int64) error
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
	err := r.coll.InsertReturning(&commentsDB)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("CommentrepositoryCreate: %w", err)
	}
	return r.mapCommentDbModelToDomain(commentsDB), nil
}

func (r commentsRepository) GetComment(id int64) (domain.Comment, error) {
	var comment comments

	err := r.coll.Find(db.Cond{"id": id}).One(&comment)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("CommentRepository GetComment: %w", err)
	}
	return r.mapCommentDbModelToDomain(comment), nil
}

func (r commentsRepository) UpdateComment(comment domain.Comment) (domain.Comment, error) {
	updateComment := r.mapCommentDBModel(comment)

	err := r.coll.Find(db.Cond{"id": updateComment.ID}).Update(&updateComment)
	if err != nil {
		return domain.Comment{}, fmt.Errorf("CommentRepository UpdateComment: %w", err)
	}

	err = r.coll.Find(db.Cond{"id": updateComment.ID}).One(&updateComment)
	if err != nil {
		return domain.Comment{}, err
	}
	return r.mapCommentDbModelToDomain(updateComment), nil
}

func (r commentsRepository) DeleteComment(id int64) error {
	var comment comments

	err := r.coll.Find(db.Cond{"id": id}).One(&comment)
	if err != nil {
		return fmt.Errorf("CommentRepository Delete: %w", err)
	}
	err = r.coll.Find(db.Cond{"id": id}).Delete()
	if err != nil {
		return fmt.Errorf("CommentRepository Delete: %w", err)
	}
	return nil
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

func (r commentsRepository) mapCommentDbModelToDomain(mcomment comments) domain.Comment {
	return domain.Comment{
		ID:     mcomment.ID,
		PostID: mcomment.PostID,
		Name:   mcomment.Name,
		Email:  mcomment.Email,
		Body:   mcomment.Body,
	}
}
