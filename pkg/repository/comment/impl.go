package comment

import (
	"context"
	"fin-asg/config/postgres"
	"fin-asg/pkg/domain/comment"
	"log"

	"gorm.io/gorm/clause"
)

type CommentRepoImpl struct {
	pgCln postgres.PostgresClient
}

func (c *CommentRepoImpl) InsertComment(ctx context.Context, input *comment.Comment) (result comment.Comment, err error) {
	log.Printf("%T - InsertComment is invoked\n", c)
	defer log.Printf("%T - InsertComment executed\n", c)

	db := c.pgCln.GetClient()

	err = db.Model(&result).Create(&input).Error

	if err != nil {
		log.Printf("error when creating comment for photo id %v\n", input.PhotoID)
	}

	result = *input

	return result, err
}

func (c *CommentRepoImpl) GetComments(ctx context.Context, userId uint64) (result []comment.Comment, err error) {
	log.Printf("%T - GetComments is invoked\n", c)
	defer log.Printf("%T - GetComments executed\n", c)

	db := c.pgCln.GetClient()

	err = db.Model(&comment.Comment{}).Where("user_id = ?", userId).Find(&result).Error

	if err != nil {
		log.Printf("error when getting photos by user id %v\n", userId)
	}

	return result, err
}

func (c *CommentRepoImpl) GetCommentById(ctx context.Context, commentId uint64) (result comment.Comment, err error) {
	log.Printf("%T - GetCommentById is invoked\n", c)
	defer log.Printf("%T - GetCommentById executed\n", c)

	db := c.pgCln.GetClient()

	err = db.Table("comments").Where("id = ?", commentId).Select("id", "message", "user_id", "photo_id").Find(&result).Error

	if err != nil {
		log.Printf("error when getting comment by id %v\n", commentId)
	}

	return result, err
}

func (c *CommentRepoImpl) UpdateComment(ctx context.Context, inputMessage string) (result comment.Comment, err error) {
	log.Printf("%T - UpdateComment is invoked\n", c)
	defer log.Printf("%T - UpdateComment executed\n", c)

	id := ctx.Value("id").(uint64)

	db := c.pgCln.GetClient()

	err = db.Model(&result).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}, {Name: "photo_id"}, {Name: "message"}, {Name: "user_id"}, {Name: "updated_at"}}}).Where("id = ?", id).Updates(comment.Comment{Message: inputMessage}).Error

	if err != nil {
		log.Printf("error when updating comment by id %v\n", id)
	}

	return result, err
}

func (c *CommentRepoImpl) DeleteComment(ctx context.Context, id uint64) (err error) {
	log.Printf("%T - DeleteComment is invoked\n", c)
	defer log.Printf("%T - DeleteComment executed\n", c)

	db := c.pgCln.GetClient()

	err = db.Where("id = ?", id).Delete(&comment.Comment{}).Error

	if err != nil {
		log.Printf("error when deleting comment by id %v \n", id)
	}

	return err
}

func NewCommentRepo(pgCln postgres.PostgresClient) comment.CommentRepo {
	return &CommentRepoImpl{pgCln: pgCln}
}
