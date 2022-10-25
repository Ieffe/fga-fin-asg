package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"fin-asg/pkg/domain/user"
	"fin-asg/pkg/usecase/crypto"

	"github.com/gin-gonic/gin"
)




type AuthMiddleware interface {
	CheckJWTAuth(ctx *gin.Context)
}

type AuthMiddlewareImpl struct {
	userUsecase user.UserUseCase
}

func (m *AuthMiddlewareImpl) CheckJWTAuth(ctx *gin.Context) {
	headerToken := ctx.Request.Header.Get("Authorization")
	bearer := strings.HasPrefix(headerToken, "Bearer")

	if !bearer {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": 97,
			"type": "UNAUTHENTICATED",
		})
		return
	}

	stringToken := strings.Split(headerToken, " ")[1]

	claim, err := crypto.VerifyJWT(ctx, stringToken)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": 97,
			"type": "UNAUTHENTICATED",
			"invalid_arg": gin.H{
				"error_type": "INVALID_TOKEN",
			},
		})
		return
	}

	userId, _ := strconv.ParseUint(claim.Subject, 0, 64)

	_, errMsg := m.userUsecase.GetUserByIdSvc(ctx, userId)

	if errMsg.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": 97,
			"type": "UNAUTHENTICATED",
			"invalid_arg": gin.H{
				"error_type":    "USER_NOT_FOUND",
				"error_message": errMsg.Error.Error(),
			},
		})
		return
	}

	// validate claim
	if claim.Issuer != "mygram.com" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": 97,
			"type": "UNAUTHENTICATED",
		})
		return
	}

	if claim.Audience != "user.mygram.com" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": 97,
			"type": "UNAUTHENTICATED",
		})
		return
	}

	if claim.Scope != "create update delete read" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": 97,
			"type": "UNAUTHENTICATED",
		})
		return
	}

	if !time.Unix(claim.NotValidBefore, 0).Before(time.Now()) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": 97,
			"type": "UNAUTHENTICATED",
		})
		return
	}

	if time.Unix(claim.ExpiredAt, 0).Before(time.Now()) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": 97,
			"type": "UNAUTHENTICATED",
		})
		return
	}

	ctx.Set("user", claim.Subject)
	ctx.Next()
}

func NewAuthMiddleware(userUsecase user.UserUseCase) AuthMiddleware {
	return &AuthMiddlewareImpl{userUsecase: userUsecase}
}
