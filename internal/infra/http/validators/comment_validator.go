package validators

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"log"
	"net/http"
	"trainee/internal/domain"
)

type CommentValidator struct {
	validator *validator.Validate
}

func NewCoordinateValidator() CommentValidator {
	return CommentValidator{
		validator: validator.New(),
	}
}

func (v CommentValidator) ValidateAndMap(request *http.Request) (domain.Comment, error) {
	var commentReq commentRequest
	err := json.NewDecoder(request.Body).Decode(&commentReq)
	if err != nil {
		log.Print(err)
		return domain.Comment{}, err
	}
	err = v.validator.Struct(commentReq)
	if err != nil {
		log.Print(err)
		return domain.Comment{}, err
	}
	return mapCommentRequestDomain(commentReq), nil
}
