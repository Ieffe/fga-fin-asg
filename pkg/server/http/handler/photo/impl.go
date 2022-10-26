package v1

import (
	"errors"
	"fin-asg/pkg/domain/message"
	"fin-asg/pkg/domain/photo"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PhotoHdlImpl struct {
	photoUseCase photo.PhotoUseCase
}

func (p *PhotoHdlImpl) InsertPhotoHdl(ctx *gin.Context) {
	log.Printf("%T - InsertPhotoHdl is invoked\n", p)
	defer log.Printf("%T - InsertPhotoHdl executed\n", p)

	var input photo.Photo

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

	log.Println("calling create photo usecase service")
	result, errMsg := p.photoUseCase.CreatePhotoSvc(ctx, input)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    01,
		"message": "photo has successfully created",
		"type":    "ACCEPTED",
		"data": gin.H{
			"id":         result.ID,
			"title":      result.Title,
			"caption":    result.Caption,
			"url":        result.Url,
			"created_at": result.CreatedAt,
			"user_id":    result.UserID,
		},
	})
}

func (p *PhotoHdlImpl) GetPhotosHdl(ctx *gin.Context) {
	log.Printf("%T - GetPhotosByUserIdHdl is invoked\n", p)
	defer log.Printf("%T - GetPhotosByUserIdHdl executed\n", p)

	stringPhotoId, isPhotoIdExist := ctx.GetQuery("id")

	if isPhotoIdExist {
		photoId, _ := strconv.ParseUint(stringPhotoId, 0, 64)
		log.Println("calling get photos by id usecase service")
		result, errMsg := p.photoUseCase.GetPhotoByIdSvc(ctx, photoId)

		if errMsg.Error != nil {
			message.ErrorResponseSwitcher(ctx, errMsg)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":    00,
			"message": fmt.Sprintf("photos id %v is found", photoId),
			"type":    "SUCCESS",
			"data":    result,
		})
		return
	}

	if isPhotoIdExist && stringPhotoId == "" {
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

	stringUserId, isUserIdExist := ctx.GetQuery("user_id")

	if isUserIdExist && stringUserId == "" {
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

	if !isUserIdExist {
		stringUserId = ctx.Value("user").(string)
	}

	userId, _ := strconv.ParseUint(stringUserId, 0, 64)

	log.Println("calling get photos by user id usecase service")
	result, errMsg := p.photoUseCase.GetPhotosByUserIdSvc(ctx, userId)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    00,
		"message": fmt.Sprintf("photos by user id %v is found", userId),
		"type":    "SUCCESS",
		"data":    result,
	})
}

func (p *PhotoHdlImpl) UpdatePhotoHdl(ctx *gin.Context) {
	log.Printf("%T - UpdatePhotoHdl is invoked\n", p)
	defer log.Printf("%T - UpdatePhotoHdl executed\n", p)

	log.Println("check photoId from path parameter")
	photoIdParam := ctx.Param("photoId")

	photoId, err := strconv.ParseUint(photoIdParam, 0, 64)

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

	log.Println("calling get photo by id usecase service")
	result, errMsg := p.photoUseCase.GetPhotoByIdSvc(ctx, photoId)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	stringUserId := ctx.Value("user").(string)
	userId, _ := strconv.ParseUint(stringUserId, 0, 64)

	log.Println("verify the photo belongs to")
	if result.UserID != userId {
		message.ErrorResponseSwitcher(ctx, message.ErrorMessage{
			Type:  "INVALID_SCOPE",
			Error: errors.New("cannot update the photo"),
		})
		return
	}

	var updatedPhoto photo.Photo

	log.Println("binding body payload from request")
	if err := ctx.ShouldBindJSON(&updatedPhoto); err != nil {
		message.ErrorResponseSwitcher(ctx, message.ErrorMessage{
			Type:  "INVALID_PAYLOAD",
			Error: errors.New("failed to bind payload"),
		})
		return
	}

	ctx.Set("photoId", result.ID)

	result, errMsg = p.photoUseCase.UpdatePhotoSvc(ctx, updatedPhoto.Title, updatedPhoto.Caption, updatedPhoto.Url)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    00,
		"message": "photo has been successfully updated",
		"type":    "SUCCESS",
		"data":    result,
	})
}

func (p *PhotoHdlImpl) DeletePhotoHdl(ctx *gin.Context) {
	log.Printf("%T - DeletePhotoHdl is invoked\n", p)
	defer log.Printf("%T - DeletePhotoHdl executed\n", p)

	log.Println("check photoId from path parameter")
	photoIdParam := ctx.Param("photoId")

	photoId, err := strconv.ParseUint(photoIdParam, 0, 64)

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

	log.Println("calling get photo by id usecase service")
	result, errMsg := p.photoUseCase.GetPhotoByIdSvc(ctx, photoId)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	stringUserId := ctx.Value("user").(string)
	userId, _ := strconv.ParseUint(stringUserId, 0, 64)

	log.Println("verify the photo belongs to")
	if result.UserID != userId {
		message.ErrorResponseSwitcher(ctx, message.ErrorMessage{
			Type:  "INVALID_SCOPE",
			Error: errors.New("cannot delete the photo"),
		})
		return
	}

	log.Println("calling delete photo usecase service")
	errMsg = p.photoUseCase.DeletePhotoSvc(ctx, photoId)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    00,
		"message": "photo has been successfully deleted",
		"type":    "SUCCESS",
	})
}

func NewPhotoHandler(photoUseCase photo.PhotoUseCase) photo.PhotoHandler {
	return &PhotoHdlImpl{photoUseCase: photoUseCase}
}
