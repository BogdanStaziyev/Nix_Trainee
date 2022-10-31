package handlers

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"trainee/internal/app"
	"trainee/internal/domain"
	"trainee/internal/infra/requests"
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
// @Param			input body domain.Post true "comment info"
// @Success 		201 {object} domain.Post
// @Failure			400 {object} error
// @Failure 		422 {object} error
// @Failure 		500 {object} error
// @Security        ApiKeyAuth
// @Router			/posts/save [post]
func (p PostHandler) SavePost(ctx echo.Context) error {
	var postRequest requests.PostRequest
	err := ctx.Bind(&postRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not decode post data"))
	}
	err = ctx.Validate(&postRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*app.JwtAccessClaim)
	id := claims.ID
	post := domain.Post{
		Title:  postRequest.Title,
		Body:   postRequest.Body,
		UserID: id,
	}
	post, err = p.service.SavePost(post)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusCreated, post)
}

// GetPost  		godoc
// @Summary 		Get Post
// @Description 	Get Post
// @Tags			Posts Actions
// @Produce 		json
// @Param			id path int true "ID"
// @Success 		200 {object} domain.Post
// @Failure 		404 {object} error
// @Security        ApiKeyAuth
// @Router			/posts/post/{id} [get]
func (p PostHandler) GetPost(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse post ID"))
	}
	post, err := p.service.GetPost(id)
	if err != nil {
		log.Print(err)
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return echo.NewHTTPError(http.StatusNotFound, err)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusOK, post)
}

// UpdatePost  		godoc
// @Summary 		Update Post
// @Description 	Update Post
// @Tags			Posts Actions
// @Accept 			json
// @Produce 		json
// @Param			input body domain.Post true "post info"
// @Success 		200 {object} domain.Post
// @Failure 		404 {object} error
// @Security        ApiKeyAuth
// @Router			/posts/update [put]
func (p PostHandler) UpdatePost(ctx echo.Context) error {
	var postRequest requests.PostRequest
	err := ctx.Bind(&postRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not decode post data")
	}
	err = ctx.Validate(&postRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse post ID"))
	}
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*app.JwtAccessClaim)
	userID := claims.ID
	post := domain.Post{
		Title:  postRequest.Title,
		Body:   postRequest.Body,
		UserID: userID,
		ID:     id,
	}
	post, err = p.service.UpdatePost(post)
	if err != nil {
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return echo.NewHTTPError(http.StatusNotFound, err)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusOK, post)
}

// DeletePost  		godoc
// @Summary 		Delete Post
// @Description 	Delete Post
// @Tags			Posts Actions
// @Produce 		json
// @Param			id path int true "ID"
// @Success 		200
// @Failure 		404 {object} error
// @Security        ApiKeyAuth
// @Router			/posts/delete/{id} [delete]
func (p PostHandler) DeletePost(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not parse post ID")
	}
	err = p.service.DeletePost(id)
	if err != nil {
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return echo.NewHTTPError(http.StatusNotFound, err)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	return ctx.NoContent(http.StatusOK)
}
