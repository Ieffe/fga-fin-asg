package photo

import (
	"context"
	"errors"
	"fin-asg/pkg/domain/message"
	"fin-asg/pkg/domain/photo"
	"fin-asg/pkg/domain/user"
	"fmt"
	"log"
	"strconv"

	"github.com/asaskevich/govalidator"
)

type PhotoUsecaseImpl struct {
	photoRepo   photo.PhotoRepo
	userUsecase user.UserUseCase
}

func (p *PhotoUsecaseImpl) CreatePhotoSvc(ctx context.Context, photo photo.Photo) (result photo.Photo, errMsg message.ErrorMessage) {
	log.Printf("%T - CreatePhotoSvc is invoked\n", p)
	defer log.Printf("%T - CreatePhotoSvc executed\n", p)

	if isValid, err := govalidator.ValidateStruct(photo); !isValid {
		switch err.Error() {
		case "photo title is required":
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "PHOTO_TITLE_IS_EMPTY",
			}
			return result, errMsg
		case "photo url is required":
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "PHOTO_URL_IS_EMPTY",
			}
			return result, errMsg
		case "invalid url format":
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "INVALID_URL_FORMAT",
			}
			return result, errMsg
		default:
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "INVALID_FORMAT",
			}
			return result, errMsg
		}
	}

	stringUserId := ctx.Value("user").(string)

	userId, _ := strconv.ParseUint(stringUserId, 0, 64)

	photo.UserID = userId

	
	result, err := p.photoRepo.InsertPhoto(ctx, &photo)

	if err != nil {
		log.Printf("error when fetching data from database: %s\n", err.Error())
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "INTERNAL_CONNECTION_PROBLEM",
		}
		return result, errMsg
	}
	return result, errMsg
}

func (p *PhotoUsecaseImpl) GetPhotosByUserIdSvc(ctx context.Context, userId uint64) (result []photo.Photo, errMsg message.ErrorMessage) {
	log.Printf("%T - GetPhotosByUserIdSvc is invoked\n", p)
	defer log.Printf("%T - GetPhotosByUserIdSvc executed\n", p)

	log.Println("calling get user by id repo")
	checkUserId, errMsg := p.userUsecase.GetUserByIdSvc(ctx, userId)

	if checkUserId.ID <= 0 {
		return result, errMsg
	}

	log.Println("calling get photo by userid repo")
	result, err := p.photoRepo.GetPhotosByUserId(ctx, userId)

	if err != nil {
		log.Printf("error when fetching data from database: %s\n", err.Error())
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "INTERNAL_CONNECTION_PROBLEM",
		}
		return result, errMsg
	}

	return result, errMsg
}

func (p *PhotoUsecaseImpl) GetPhotoByIdSvc(ctx context.Context, id uint64) (result photo.Photo, errMsg message.ErrorMessage) {
	log.Printf("%T - GetPhotoByIdSvc is invoked\n", p)
	defer log.Printf("%T - GetPhotoByIdSvc executed\n", p)

	log.Println("calling get photo by id repo")
	result, err := p.photoRepo.GetPhotoById(ctx, id)

	if result.ID <= 0 {
		log.Printf("photo with id %v not found", id)

		err = fmt.Errorf("photo with id %v not found", id)
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "PHOTO_NOT_FOUND",
		}
		return result, errMsg
	}

	if err != nil {
		log.Printf("error when fetching data from database: %s\n", err.Error())
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "INTERNAL_CONNECTION_PROBLEM",
		}
		return result, errMsg
	}

	return result, errMsg
}

func (p *PhotoUsecaseImpl) UpdatePhotoSvc(ctx context.Context, title string, caption string, url string) (result photo.Photo, errMsg message.ErrorMessage) {
	log.Printf("%T - UpdatePhotoSvc is invoked\n", p)
	defer log.Printf("%T - UpdatePhotoSvc executed\n", p)

	if title == "" {
		errMsg := message.ErrorMessage{
			Error: errors.New("photo title is required"),
			Type:  "PHOTO_TITLE_IS_EMPTY",
		}
		return result, errMsg
	}
	if url == "" {
		errMsg := message.ErrorMessage{
			Error: errors.New("photo url is required"),
			Type:  "PHOTO_URL_IS_EMPTY",
		}
		return result, errMsg
	}

	if !govalidator.IsURL(url) {
		errMsg := message.ErrorMessage{
			Error: errors.New("invalid url format"),
			Type:  "INVALID_URL_FORMAT",
		}
		return result, errMsg
	}

	result, err := p.photoRepo.UpdatePhoto(ctx, title, caption, url)

	if err != nil {
		log.Printf("error when fetching data from database: %s\n", err.Error())
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "INTERNAL_CONNECTION_PROBLEM",
		}
		return result, errMsg
	}

	return result, errMsg
}

func (p *PhotoUsecaseImpl) DeletePhotoSvc(ctx context.Context, photoId uint64) (errMsg message.ErrorMessage) {
	log.Printf("%T - DeletePhotoSvc is invoked\n", p)
	defer log.Printf("%T - DeletePhotoSvc executed\n", p)

	log.Println("calling delete photo repo")
	err := p.photoRepo.DeletePhoto(ctx, photoId)

	if err != nil {
		log.Printf("error when fetching data from database: %s\n", err.Error())
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "INTERNAL_CONNECTION_PROBLEM",
		}
		return errMsg
	}

	return errMsg
}

func NewPhotoUsecase(photoRepo photo.PhotoRepo, userUsecase user.UserUseCase) photo.PhotoUseCase {
	return &PhotoUsecaseImpl{photoRepo: photoRepo, userUsecase: userUsecase}
}
