package validators

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"log"
	"net/http"
	"trainee/internal/domain"
)

type PostValidator struct {
	validator *validator.Validate
}

func NewPostValidator() PostValidator {
	return PostValidator{
		validator: validator.New(),
	}
}

func (v PostValidator) ValidateAndMap(request *http.Request) (domain.Post, error) {
	var postReq postRequest
	err := json.NewDecoder(request.Body).Decode(&postReq)
	if err != nil {
		log.Print(err)
		return domain.Post{}, err
	}
	err = v.validator.Struct(postReq)
	if err != nil {
		log.Print(err)
		return domain.Post{}, err
	}
	return mapPostRequestDomain(postReq), nil
}
