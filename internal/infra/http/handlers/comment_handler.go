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

type CommentHandler struct {
	service    app.CommentService
	usrService app.UserService
}

func NewCommentHandler(s app.CommentService, u app.UserService) CommentHandler {
	return CommentHandler{
		service:    s,
		usrService: u,
	}
}

// SaveComment 		godoc
// @Summary 		Save Comment
// @Description 	Save Comment
// @Tags			Comments Actions
// @Accept 			json
// @Produce 		json
// @Param			post_id path int true "PostID"
// @Param			input body requests.CommentRequest true "comment info"
// @Success 		201 {object} response.CommentResponse
// @Failure			400 {object} response.Error
// @Failure 		422 {object} response.Error
// @Failure 		500 {object} response.Error
// @Security        ApiKeyAuth
// @Router			/api/v1/comments/save/{post_id} [post]
func (c CommentHandler) SaveComment(ctx echo.Context) error {
	var commentRequest requests.CommentRequest
	err := ctx.Bind(&commentRequest)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "could not decode comment data")
	}
	err = ctx.Validate(&commentRequest)
	if err != nil {
		log.Print(err)
		return response.ErrorResponse(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	postID, err := strconv.ParseInt(ctx.Param("post_id"), 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "could not parse comment ID")
	}
	jwtUser := ctx.Get("user").(*jwt.Token)
	claims := jwtUser.Claims.(*app.JwtAccessClaim)
	user, err := c.usrService.FindByID(claims.ID)
	if err != nil {
		log.Printf("SaveComment error, %s", err)
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return response.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		} else {
			return response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		}
	}
	comment := domain.Comment{
		PostID: postID,
		Name:   user.Name,
		Email:  user.Email,
		Body:   commentRequest.Body,
	}
	comment, err = c.service.SaveComment(comment)
	if err != nil {
		log.Print(err)
		return response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	commentResponse := domain.Comment.DomainToResponse(comment)
	return response.Response(ctx, http.StatusCreated, commentResponse)
}

// GetComment 		godoc
// @Summary 		Get Comment
// @Description 	Get Comment
// @Tags			Comments Actions
// @Produce 		json
// @Param			id path int true "ID"
// @Success 		200 {object} response.CommentResponse
// @Failure			400 {object} response.Error
// @Failure			404 {object} response.Error
// @Failure			500 {object} response.Error
// @Security 		ApiKeyAuth
// @Router			/api/v1/comments/comment/{id} [get]
func (c CommentHandler) GetComment(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, "could not parse comment ID")
	}
	comment, err := c.service.GetComment(id)
	if err != nil {
		log.Print("commentService error", err)
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return response.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		} else {
			return response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		}
	}
	commentResponse := domain.Comment.DomainToResponse(comment)
	return response.Response(ctx, http.StatusOK, commentResponse)
}

// UpdateComment 	godoc
// @Summary 		Update Comment
// @Description 	Update Comment
// @Tags			Comments Actions
// @Accept 			json
// @Produce 		json
// @Param			id path int true "ID"
// @Param			input body requests.CommentRequest true "comment info"
// @Success 		200 {object} response.CommentResponse
// @Failure			400 {object} response.Error
// @Failure			422 {object} response.Error
// @Failure			404 {object} response.Error
// @Security 		ApiKeyAuth
// @Router			/api/v1/comments/update/{id} [put]
func (c CommentHandler) UpdateComment(ctx echo.Context) error {
	var commentRequest requests.CommentRequest
	err := ctx.Bind(&commentRequest)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprint(err.Error()+"could not decode comment data"))
	}
	err = ctx.Validate(&commentRequest)
	if err != nil {
		log.Print(err)
		return response.ErrorResponse(ctx, http.StatusUnprocessableEntity, err.Error())
	}
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprint(err.Error()+"could not parse comment ID"))
	}
	comment, err := c.service.UpdateComment(commentRequest.Body, id)
	if err != nil {
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return response.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		} else {
			return response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		}
	}
	commentResponse := domain.Comment.DomainToResponse(comment)
	return response.Response(ctx, http.StatusOK, commentResponse)
}

// DeleteComment	godoc
// @Summary 		Delete Comment
// @Description 	Delete Comment
// @Tags			Comments Actions
// @Param			id path int true "ID"
// @Success 		200 {object} response.Data
// @Failure			400	{object} response.Error
// @Failure			404 {object} response.Error
// @Failure			500 {object} response.Error
// @Security 		ApiKeyAuth
// @Router			/api/v1/comments/delete/{id} [delete]
func (c CommentHandler) DeleteComment(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error()+"could not parse comment ID")
	}
	err = c.service.DeleteComment(id)
	if err != nil {
		log.Print(err)
		if strings.HasSuffix(err.Error(), "upper: no more rows in this result set") {
			return response.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		} else {
			return response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		}
	}
	return response.MessageResponse(ctx, http.StatusOK, "Comment Delete")
}
