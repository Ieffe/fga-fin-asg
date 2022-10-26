package v1

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"fin-asg/pkg/domain/message"
	"fin-asg/pkg/domain/social"

	"github.com/gin-gonic/gin"
)

type SocialHdlImpl struct {
	socialUseCase social.SocialUseCase
}

func (c *SocialHdlImpl) InsertSocialHdl(ctx *gin.Context) {
	log.Printf("%T - InsertSocialHdl is invoked\n", c)
	defer log.Printf("%T - InsertSocialHdl executed\n", c)

	var inputSocial social.Social

	log.Println("binding body payload from request")
	if err := ctx.ShouldBindJSON(&inputSocial); err != nil {
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

	log.Println("calling create social usecase service")
	result, errMsg := c.socialUseCase.InsertSocialSvc(ctx, inputSocial)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    01,
		"message": "social  has successfully created",
		"type":    "ACCEPTED",
		"data": gin.H{
			"id":         result.ID,
			"name":       result.Name,
			"url":        result.URL,
			"created_at": result.CreatedAt,
			"user_id":    result.UserID,
		},
	})
}

func (c *SocialHdlImpl) GetSocialsHdl(ctx *gin.Context) {
	log.Printf("%T - GetSocialsIdHdl is invoked\n", c)
	defer log.Printf("%T - GetSocialsHdl executed\n", c)

	stringUserId := ctx.Value("user").(string)

	userId, _ := strconv.ParseUint(stringUserId, 0, 64)

	log.Println("calling get socials usecase service")
	result, errMsg := c.socialUseCase.GetSocialsSvc(ctx)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    00,
		"message": fmt.Sprintf("social s by user id %v is found", userId),
		"type":    "SUCCESS",
		"data":    result,
	})
}

func (c *SocialHdlImpl) UpdateSocialHdl(ctx *gin.Context) {
	log.Printf("%T - UpdateSocialHdl is invoked\n", c)
	defer log.Printf("%T - UpdateSocialHdl executed\n", c)

	log.Println("check social  id from path parameter")
	idParam := ctx.Param("socialId")

	id, err := strconv.ParseUint(idParam, 0, 64)

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

	log.Println("calling get social  by id usecase service")
	result, errMsg := c.socialUseCase.GetSocialByIdSvc(ctx, id)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	stringUserId := ctx.Value("user").(string)
	userId, _ := strconv.ParseUint(stringUserId, 0, 64)

	log.Println("verify the social  belongs to")
	if result.UserID != userId {
		message.ErrorResponseSwitcher(ctx, message.ErrorMessage{
			Type:  "INVALID_SCOPE",
			Error: errors.New("cannot update the social "),
		})
		return
	}

	var updatedSocial social.Social

	log.Println("binding body payload from request")
	if err := ctx.ShouldBindJSON(&updatedSocial); err != nil {
		message.ErrorResponseSwitcher(ctx, message.ErrorMessage{
			Type:  "INVALID_PAYLOAD",
			Error: errors.New("failed to bind payload"),
		})
		return
	}

	ctx.Set("id", id)

	updateResult, errMsg := c.socialUseCase.UpdateSocialSvc(ctx, updatedSocial)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    01,
		"message": "social  has been successfully updated",
		"type":    "ACCEPTED",
		"data":    updateResult,
	})
}

func (c *SocialHdlImpl) DeleteSocialHdl(ctx *gin.Context) {
	log.Printf("%T - DeleteSocialHdl is invoked\n", c)
	defer log.Printf("%T - DeleteSocialHdl executed\n", c)

	log.Println("check id from path parameter")
	idParam := ctx.Param("socialId")

	id, err := strconv.ParseUint(idParam, 0, 64)

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

	log.Println("calling get social  by id usecase service")
	result, errMsg := c.socialUseCase.GetSocialByIdSvc(ctx, id)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	stringUserId := ctx.Value("user").(string)
	userId, _ := strconv.ParseUint(stringUserId, 0, 64)

	log.Println("verify the social  belongs to")
	if result.UserID != userId {
		message.ErrorResponseSwitcher(ctx, message.ErrorMessage{
			Type:  "INVALID_SCOPE",
			Error: errors.New("cannot delete the social "),
		})
		return
	}

	log.Println("calling delete social  usecase service")
	errMsg = c.socialUseCase.DeleteSocialSvc(ctx, id)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    01,
		"message": "social  has been successfully deleted",
		"type":    "ACCEPTED",
	})
}

func NewSocialHandler(socialUseCase social.SocialUseCase) social.SocialHandler {
	return &SocialHdlImpl{socialUseCase: socialUseCase}
}
