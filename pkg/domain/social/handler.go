package social

import "github.com/gin-gonic/gin"

type SocialHandler interface {
	InsertSocialHdl(ctx *gin.Context)
	GetSocialsHdl(ctx *gin.Context)
	UpdateSocialHdl(ctx *gin.Context)
	DeleteSocialHdl(ctx *gin.Context)
}