package v1

import (
	"errors"
	"fin-asg/pkg/domain/comment"
	"fin-asg/pkg/domain/message"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHdlImpl struct {
	commentUsecase comment.CommentUsecase
}

func (c *CommentHdlImpl) InsertCommentHdl(ctx *gin.Context) {
	log.Printf("%T - InsertCommentHdl is invoked\n", c)
	defer log.Printf("%T - InsertCommentHdl executed\n", c)

	var input comment.Comment

	log.Println("binding body payload from request")
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    96,
			"type":    "BAD_REQUEST",
			"message": "Failed to bind payload",
			"invalid_arg": gin.H{
				"error_type":    "INVALID_PAYLOAD",
				"error_message": err.Error(),
			},
		})
		return
	}

	log.Println("calling insert comment usecase service")
	result, errMsg := c.commentUsecase.InsertCommentSvc(ctx, input)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    01,
		"message": "comment has successfully created",
		"type":    "ACCEPTED",
		"data": gin.H{
			"id":         result.ID,
			"message":    result.Message,
			"photo_id":   result.PhotoID,
			"user_id":    result.UserID,
			"created_at": result.CreatedAt,
		},
	})
}

func (c *CommentHdlImpl) GetCommentsHdl(ctx *gin.Context) {
	log.Printf("%T - GetCommentsIdHdl is invoked\n", c)
	defer log.Printf("%T - GetCommentsHdl executed\n", c)

	stringUserId := ctx.Value("user").(string)

	userId, _ := strconv.ParseUint(stringUserId, 0, 64)

	log.Println("calling get comments by user id usecase service")
	result, errMsg := c.commentUsecase.GetCommentsSvc(ctx)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    00,
		"message": fmt.Sprintf("comments by user id %v is found", userId),
		"type":    "SUCCESS",
		"data":    result,
	})
}

func (c *CommentHdlImpl) UpdateCommentHdl(ctx *gin.Context) {
	log.Printf("%T - UpdateCommentHdl is invoked\n", c)
	defer log.Printf("%T - UpdateCommentHdl executed\n", c)

	log.Println("check commentId from path parameter")
	commentIdParam := ctx.Param("commentId")

	commentId, err := strconv.ParseUint(commentIdParam, 0, 64)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    96,
			"type":    "BAD_REQUEST",
			"message": "invalid params",
			"invalid_arg": gin.H{
				"error_type":    "INVALID_PARAMS",
				"error_message": "invalid params",
			},
		})
		return
	}

	log.Println("calling get comment by id usecase service")
	result, errMsg := c.commentUsecase.GetCommentByIdSvc(ctx, commentId)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	stringUserId := ctx.Value("user").(string)
	userId, _ := strconv.ParseUint(stringUserId, 0, 64)

	log.Println("verify the comment belongs to")
	if result.UserID != userId {
		message.ErrorResponseSwitcher(ctx, message.ErrorMessage{
			Type:  "INVALID_SCOPE",
			Error: errors.New("cannot update the comment"),
		})
		return
	}

	var updatedComment comment.Comment

	log.Println("binding body payload from request")
	if err := ctx.ShouldBindJSON(&updatedComment); err != nil {
		message.ErrorResponseSwitcher(ctx, message.ErrorMessage{
			Type:  "INVALID_PAYLOAD",
			Error: errors.New("failed to bind payload"),
		})
		return
	}

	ctx.Set("commentId", commentId)

	updateResult, errMsg := c.commentUsecase.UpdateCommentSvc(ctx, updatedComment.Message)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    01,
		"message": "comment has been successfully updated",
		"type":    "ACCEPTED",
		"data":    updateResult,
	})
}

func (c *CommentHdlImpl) DeleteCommentHdl(ctx *gin.Context) {
	log.Printf("%T - DeleteCommentHdl is invoked\n", c)
	defer log.Printf("%T - DeleteCommentHdl executed\n", c)

	log.Println("check commentId from path parameter")
	commentIdParam := ctx.Param("commentId")

	commentId, err := strconv.ParseUint(commentIdParam, 0, 64)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    96,
			"type":    "BAD_REQUEST",
			"message": "invalid params",
			"invalid_arg": gin.H{
				"error_type":    "INVALID_PARAMS",
				"error_message": "invalid params",
			},
		})
		return
	}

	log.Println("calling get comment by id usecase service")
	result, errMsg := c.commentUsecase.GetCommentByIdSvc(ctx, commentId)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	stringUserId := ctx.Value("user").(string)
	userId, _ := strconv.ParseUint(stringUserId, 0, 64)

	log.Println("verify the comment belongs to")
	if result.UserID != userId {
		message.ErrorResponseSwitcher(ctx, message.ErrorMessage{
			Type:  "INVALID_SCOPE",
			Error: errors.New("cannot delete the comment"),
		})
		return
	}

	log.Println("calling delete comment usecase service")
	errMsg = c.commentUsecase.DeleteCommentSvc(ctx, commentId)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    01,
		"message": "comment has been successfully deleted",
		"type":    "ACCEPTED",
	})
}

func NewCommentHandler(commentUsecase comment.CommentUsecase) comment.CommentHandler {
	return &CommentHdlImpl{commentUsecase: commentUsecase}
}
