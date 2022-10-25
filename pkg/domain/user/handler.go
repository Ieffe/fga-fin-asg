package user
import (
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetUserByIdHdl(ctx *gin.Context)
	InsertUserHdl(ctx *gin.Context)
	UpdateUserHdl(ctx *gin.Context)
	DeleteUserHdl(ctx *gin.Context)
}

