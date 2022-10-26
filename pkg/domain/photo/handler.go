package photo

import "github.com/gin-gonic/gin"

type PhotoHandler interface {
	InsertPhotoHdl(ctx *gin.Context)
	GetPhotosHdl(ctx *gin.Context)
	UpdatePhotoHdl(ctx *gin.Context)
	DeletePhotoHdl(ctx *gin.Context)
}