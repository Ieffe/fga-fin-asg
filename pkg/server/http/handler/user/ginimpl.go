package user

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"fin-asg/pkg/domain/message"
	"fin-asg/pkg/domain/user"
	"fin-asg/pkg/helper"

	"github.com/gin-gonic/gin"
)

type UserHdlImpl struct {
	UserUseCase user.UserUseCase
}

func NewUserHandler(UserUseCase user.UserUseCase) user.UserHandler {
	return &UserHdlImpl{UserUseCase: UserUseCase}
}

func (u *UserHdlImpl) GetUserByIdHdl(ctx *gin.Context) {
	log.Printf("%T - GetUserByIdHdl is invoked\n", u)
	defer log.Printf("%T - GetUserByIdHdl executed\n", u)

	log.Println("binding body payload from request")
	var user user.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
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

	currentAge := helper.ConvertDOBToCurrentAge(time.Time(user.Dob), time.Now())

	log.Println("checking user age")
	if currentAge < 8 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    96,
			"message": "Your minimal age should be 8 years old",
			"type":    "BAD_REQUEST",
		})
		return
	}

	log.Println("calling register user service usecase")
	result, errMsg := u.UserUseCase.InsertUserSvc(ctx, user)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    01,
		"message": "user has successfully registered",
		"type":    "ACCEPTED",
		"data": gin.H{
			"age":      currentAge,
			"email":    result.Email,
			"id":       result.ID,
			"username": result.Username,
		},
	})
}

func (u *UserHdlImpl) InsertUserHdl(ctx *gin.Context) {
	log.Printf("%T - RegisterUserHdl is invoked\n", u)
	defer log.Printf("%T - RegisterUserHdl executed\n", u)

	log.Println("binding body payload from request")
	var user user.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
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

	currentAge := helper.ConvertDOBToCurrentAge(time.Time(user.Dob), time.Now())

	log.Println("calling register user service usecase")
	result, errMsg := u.UserUseCase.InsertUserSvc(ctx, user)

	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    01,
		"message": "user has successfully registered",
		"type":    "ACCEPTED",
		"data": gin.H{
			"age":      currentAge,
			"email":    result.Email,
			"id":       result.ID,
			"username": result.Username,
		},
	})
}

func (u *UserHdlImpl) UpdateUserHdl(ctx *gin.Context) {
	log.Printf("%T - GetUserByIdHdl is invoked\n", u)
	defer log.Printf("%T - GetUserByIdHdl executed\n", u)

	log.Println("check id from path parameter")
	userIdParam := ctx.Param("id")

	userId, err := strconv.ParseUint(userIdParam, 0, 64)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": 96,
			"type": "BAD_REQUEST",
			"invalid_arg": gin.H{
				"error_type":    "INVALID_FORMAT",
				"error_message": err.Error(),
			},
		})
		return
	}

	log.Println("calling get user by id service usecase")
	result, errMsg := u.UserUseCase.GetUserByIdSvc(ctx, userId)
	if errMsg.Error != nil {
		message.ErrorResponseSwitcher(ctx, errMsg)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    00,
		"message": "user is found",
		"type":    "SUCCESS",
		"data": gin.H{
			"id":            result.ID,
			"username":      result.Username,
			"social_medias": result.Social,
		},
	})
}

func (u *UserHdlImpl) DeleteUserHdl(ctx *gin.Context) {
	
}


