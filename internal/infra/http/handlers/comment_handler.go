package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"trainee/internal/app"
	"trainee/internal/domain"
)

type CommentHandler struct {
	service app.CommentService
}

func NewCommentHandler(s app.CommentService) CommentHandler {
	return CommentHandler{
		service: s,
	}
}

// SaveComment 		godoc
// @Summary 		Save Comment
// @Description 	Save Comment
// @Tags			Comments Actions
// @Accept 			json
// @Produce 		json
// @Param			input body domain.Comment true "comment info"
// @Success 		201 {object} domain.Comment
// @Failure			400 {object} error
// @Failure 		422 {object} error
// @Failure 		500 {object} error
// @Security        ApiKeyAuth
// @Router			/comments/save [post]
func (c CommentHandler) SaveComment(ctx echo.Context) error {
	var comment domain.Comment
	err := ctx.Bind(&comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not decode comment data"))
	}
	err = ctx.Validate(&comment)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	comment, err = c.service.SaveComment(comment)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusCreated, comment)
}

// GetComment 		godoc
// @Summary 		Get Comment
// @Description 	Get Comment
// @Tags			Comments Actions
// @Produce 		json
// @Param			id path int true "ID"
// @Success 		200 {object} domain.Comment
// @Failure			404 {object} error
// @Security 		ApiKeyAuth
// @Router			/comments/comment/{id} [get]
func (c CommentHandler) GetComment(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse comment ID"))
	}
	comment, err := c.service.GetComment(id)
	if err != nil {
		log.Print("commentService error", err)
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return echo.NewHTTPError(http.StatusNotFound, err)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusOK, comment)
}

// UpdateComment 	godoc
// @Summary 		Update Comment
// @Description 	Update Comment
// @Tags			Comments Actions
// @Accept 			json
// @Produce 		json
// @Param			input body domain.Comment true "comment info"
// @Success 		200 {object} domain.Comment
// @Failure			422 {object} error
// @Failure			404 {object} error
// @Security 		ApiKeyAuth
// @Router			/comments/update [put]
func (c CommentHandler) UpdateComment(ctx echo.Context) error {
	var comment domain.Comment
	err := ctx.Bind(&comment)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not decode comment data"))
	}
	err = ctx.Validate(&comment)
	if err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	comment, err = c.service.UpdateComment(comment)
	if err != nil {
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return echo.NewHTTPError(http.StatusNotFound, err)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusOK, comment)
}

// DeleteComment	godoc
// @Summary 		Delete Comment
// @Description 	Delete Comment
// @Tags			Comments Actions
// @Param			id path int true "ID"
// @Success 		200
// @Failure			404 {object} error
// @Security 		ApiKeyAuth
// @Router			/comments/delete/{id} [delete]
func (c CommentHandler) DeleteComment(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse comment ID"))
	}
	err = c.service.DeleteComment(id)
	if err != nil {
		log.Print(err)
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return echo.NewHTTPError(http.StatusNotFound, err)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	return ctx.NoContent(http.StatusOK)
}
