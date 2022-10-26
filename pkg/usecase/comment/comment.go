package comment

import (
	"context"
	"errors"
	"fin-asg/pkg/domain/comment"
	"fin-asg/pkg/domain/message"
	"fin-asg/pkg/domain/photo"
	"fmt"
	"log"
	"strconv"

	"github.com/asaskevich/govalidator"
)

type CommentUsecaseImpl struct {
	commentRepo  comment.CommentRepo
	photoUsecase photo.PhotoUseCase
}

func (c *CommentUsecaseImpl) InsertCommentSvc(ctx context.Context, comment comment.Comment) (result comment.Comment, errMsg message.ErrorMessage) {
	log.Printf("%T - InsertCommentSvc is invoked\n", c)
	defer log.Printf("%T - InsertCommentSvc executed\n", c)

	if isValid, err := govalidator.ValidateStruct(comment); !isValid {
		switch err.Error() {
		case "comment message is required":
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "COMMENT_MESSAGE_IS_EMPTY",
			}
			return result, errMsg
		case "photo id is required":
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "PHOTO_ID_IS_EMPTY",
			}
			return result, errMsg
		default:
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "INVALID_FORMAT",
			}
			return result, errMsg
		}
	}

	log.Println("check photo is exist or not")
	_, errMsg = c.photoUsecase.GetPhotoByIdSvc(ctx, comment.PhotoID)

	if errMsg.Error != nil {
		return result, errMsg
	}

	stringUserId := ctx.Value("user").(string)

	userId, _ := strconv.ParseUint(stringUserId, 0, 64)

	comment.UserID = userId

	log.Println("calling insert comment repo")
	result, err := c.commentRepo.InsertComment(ctx, &comment)

	if err != nil {
		log.Printf("error when fetching data from database: %s\n", err.Error())
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "INTERNAL_CONNECTION_PROBLEM",
		}
		return result, errMsg
	}
	return result, errMsg
}

func (c *CommentUsecaseImpl) GetCommentsSvc(ctx context.Context) (result []comment.Comment, errMsg message.ErrorMessage) {
	log.Printf("%T - GetCommentsSvc is invoked\n", c)
	defer log.Printf("%T - GetCommentsSvc executed\n", c)

	stringUserId := ctx.Value("user").(string)

	userId, _ := strconv.ParseUint(stringUserId, 0, 64)

	log.Println("calling get comment by userid repo")
	result, err := c.commentRepo.GetComments(ctx, userId)

	if err != nil {
		log.Printf("error when fetching data from database: %s\n", err.Error())
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "INTERNAL_CONNECTION_PROBLEM",
		}
		return result, errMsg
	}

	return result, errMsg
}

func (c *CommentUsecaseImpl) GetCommentByIdSvc(ctx context.Context, commentId uint64) (result comment.Comment, errMsg message.ErrorMessage) {
	log.Printf("%T - GetCommentByIdSvc is invoked\n", c)
	defer log.Printf("%T - GetCommentByIdSvc executed\n", c)

	log.Println("calling get comment by id repo")
	result, err := c.commentRepo.GetCommentById(ctx, commentId)

	if result.ID <= 0 {
		log.Printf("comment with id %v not found", commentId)

		err = fmt.Errorf("comment with id %v not found", commentId)
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "COMMENT_NOT_FOUND",
		}
		return result, errMsg
	}

	if err != nil {
		log.Printf("error when fetching data from database: %s\n", err.Error())
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "INTERNAL_CONNECTION_PROBLEM",
		}
		return result, errMsg
	}

	return result, errMsg
}

func (c *CommentUsecaseImpl) UpdateCommentSvc(ctx context.Context, inputMessage string) (result comment.Comment, errMsg message.ErrorMessage) {
	log.Printf("%T - UpdateCommentSvc is invoked\n", c)
	defer log.Printf("%T - UpdateCommentSvc executed\n", c)

	if inputMessage == "" {
		errMsg := message.ErrorMessage{
			Error: errors.New("comment message is required"),
			Type:  "COMMENT_MESSAGE_IS_EMPTY",
		}
		return result, errMsg
	}

	result, err := c.commentRepo.UpdateComment(ctx, inputMessage)

	if err != nil {
		log.Printf("error when fetching data from database: %s\n", err.Error())
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "INTERNAL_CONNECTION_PROBLEM",
		}
		return result, errMsg
	}

	return result, errMsg
}

func (c *CommentUsecaseImpl) DeleteCommentSvc(ctx context.Context, commentId uint64) (errMsg message.ErrorMessage) {
	log.Printf("%T - DeleteCommentSvc is invoked\n", c)
	defer log.Printf("%T - DeleteCommentSvc executed\n", c)

	log.Println("calling delete comment repo")
	err := c.commentRepo.DeleteComment(ctx, commentId)

	if err != nil {
		log.Printf("error when fetching data from database: %s\n", err.Error())
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "INTERNAL_CONNECTION_PROBLEM",
		}
		return errMsg
	}

	return errMsg
}

func NewCommentUsecase(commentRepo comment.CommentRepo, photoUseCase photo.PhotoUseCase) comment.CommentUsecase {
	return &CommentUsecaseImpl{commentRepo: commentRepo, photoUsecase: photoUseCase}
}
