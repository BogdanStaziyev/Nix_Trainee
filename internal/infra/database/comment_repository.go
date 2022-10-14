package database

import (
	"fmt"
	"github.com/upper/db/v4"
	"trainee/internal/domain"
)

const CommentTable = "commentses"

type comments struct {
	Id     int64  `db:"id,omitempty"`
	PostID int64  `db:"post_id"`
	Name   string `db:"name"`
	Email  string `db:"email"`
	Body   string `db:"body"`
}

type CommentRepo interface {
	SaveComment(comment domain.Comment) (domain.Comment, error)
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

func (r commentsRepository) mapCommentDBModel(comment domain.Comment) comments {
	return comments{
		Id:     comment.Id,
		PostID: comment.PostId,
		Name:   comment.Name,
		Email:  comment.Email,
		Body:   comment.Body,
	}
}

func (r commentsRepository) mapCommentDbModelToDomain(mcomment comments) domain.Comment {
	return domain.Comment{
		Id:     mcomment.Id,
		PostId: mcomment.PostID,
		Name:   mcomment.Name,
		Email:  mcomment.Email,
		Body:   mcomment.Body,
	}
}
