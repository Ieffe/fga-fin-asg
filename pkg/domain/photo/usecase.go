package photo

import (
	"context"
	"fin-asg/pkg/domain/message"
)

type PhotoUseCase interface {
	CreatePhotoSvc(ctx context.Context, photo Photo) (result Photo, errMsg message.ErrorMessage)
	GetPhotosByUserIdSvc(ctx context.Context, userId uint64) (result []Photo, errMsg message.ErrorMessage)
	GetPhotoByIdSvc(ctx context.Context, photoId uint64) (result Photo, errMsg message.ErrorMessage)
	UpdatePhotoSvc(ctx context.Context, title string, caption string, url string) (result Photo, errMsg message.ErrorMessage)
	DeletePhotoSvc(ctx context.Context, photoId uint64) (errMsg message.ErrorMessage)
}