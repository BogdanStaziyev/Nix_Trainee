package handlers

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"strings"
	"trainee/internal/app"
	"trainee/internal/domain"
	"trainee/internal/infra/http/requests"
	"trainee/internal/infra/http/response"
)

type PostHandler struct {
	service app.PostService
}

func NewPostHandler(s app.PostService) PostHandler {
	return PostHandler{
		service: s,
	}
}

// SavePost 		godoc
// @Summary 		Save Post
// @Description 	Save Post
// @Tags			Posts Actions
// @Accept 			json
// @Produce 		json
// @Param			input body requests.PostRequest true "comment info"
// @Success 		201 {object} response.PostResponse
// @Failure			400 {object} response.Error
// @Failure 		422 {object} response.Error
// @Failure 		500 {object} response.Error
// @Security        ApiKeyAuth
// @Router			/api/v1/posts/save [post]
func (p PostHandler) SavePost(ctx echo.Context) error {
	var postRequest requests.PostRequest
	err := ctx.Bind(&postRequest)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not decode post data")
	}
	err = ctx.Validate(&postRequest)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Could not validate post data")
	}
	token := ctx.Get("user").(*jwt.Token)
	post, err := p.service.SavePost(postRequest, token)
	if err != nil {
		log.Print(err)
		return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Could not save new post: %s", err))
	}
	postResponse := domain.Post.DomainToResponse(post)
	return response.Response(ctx, http.StatusCreated, postResponse)
}

// GetPost  		godoc
// @Summary 		Get Post
// @Description 	Get Post
// @Tags			Posts Actions
// @Produce 		json
// @Param			id path int true "ID"
// @Success 		200 {object} response.PostResponse
// @Failure 		400 {object} response.Error
// @Failure 		404 {object} response.Error
// @Failure 		500 {object} response.Error
// @Security        ApiKeyAuth
// @Router			/api/v1/posts/post/{id} [get]
func (p PostHandler) GetPost(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not parse post ID")
	}
	post, err := p.service.GetPost(id)
	if err != nil {
		log.Print(err)
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return response.ErrorResponse(ctx, http.StatusNotFound, fmt.Sprintf("Could not get post: %s", err))
		} else {
			return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Could not get post: %s", err))
		}
	}
	postResponse := domain.Post.DomainToResponse(post)
	return response.Response(ctx, http.StatusOK, postResponse)
}

// UpdatePost  		godoc
// @Summary 		Update Post
// @Description 	Update Post
// @Tags			Posts Actions
// @Accept 			json
// @Produce 		json
// @Param			id path int true "ID"
// @Param			input body requests.PostRequest true "post info"
// @Success 		200 {object} response.PostResponse
// @Failure 		400 {object} response.Error
// @Failure 		422 {object} response.Error
// @Failure 		400 {object} response.Error
// @Failure 		404 {object} response.Error
// @Failure 		500 {object} response.Error
// @Security        ApiKeyAuth
// @Router			/api/v1/posts/update/{id} [put]
func (p PostHandler) UpdatePost(ctx echo.Context) error {
	var postRequest requests.PostRequest
	err := ctx.Bind(&postRequest)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not decode post data")
	}
	err = ctx.Validate(&postRequest)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Could not validate post data")
	}
	postID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not parse post ID")
	}
	post, err := p.service.UpdatePost(postRequest, postID)
	if err != nil {
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return response.ErrorResponse(ctx, http.StatusNotFound, fmt.Sprintf("Could not get post: %s", err))
		} else {
			return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Could not get post: %s", err))
		}
	}
	postResponse := domain.Post.DomainToResponse(post)
	return response.Response(ctx, http.StatusOK, postResponse)
}

// DeletePost  		godoc
// @Summary 		Delete Post
// @Description 	Delete Post
// @Tags			Posts Actions
// @Produce 		json
// @Param			id path int true "ID"
// @Success 		200 {object} response.Data
// @Failure 		400 {object} response.Error
// @Failure 		404 {object} response.Error
// @Failure 		500 {object} response.Error
// @Security        ApiKeyAuth
// @Router			/api/v1/posts/delete/{id} [delete]
func (p PostHandler) DeletePost(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "Could not parse post ID")
	}
	err = p.service.DeletePost(id)
	if err != nil {
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return response.ErrorResponse(ctx, http.StatusNotFound, fmt.Sprintf("Could not get post: %s", err))
		} else {
			return response.ErrorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("Could not get post: %s", err))
		}
	}
	return response.MessageResponse(ctx, http.StatusOK, "Post successfully delete")
}
