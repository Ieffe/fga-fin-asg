package photo

import "context"

type PhotoRepo interface {
	InsertPhoto(ctx context.Context, photo *Photo) (result Photo, err error)
	GetPhotosByUserId(ctx context.Context, userId uint64) (result []Photo, err error)
	GetPhotoById(ctx context.Context, id uint64) (result Photo, err error)
	UpdatePhoto(ctx context.Context, title string, caption string, url string) (result Photo, err error)
	DeletePhoto(ctx context.Context, id uint64) (err error)
}