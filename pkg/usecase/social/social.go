package social

import (
	"context"
	"fin-asg/pkg/domain/message"
	"fin-asg/pkg/domain/social"
	"fmt"
	"log"
	"strconv"

	"github.com/asaskevich/govalidator"
)

type SocialUseCaseImpl struct {
	socialRepo social.SocialRepo
}

func (s *SocialUseCaseImpl) InsertSocialSvc(ctx context.Context, social social.Social) (result social.Social, errMsg message.ErrorMessage) {
	log.Printf("%T - InsertSocialSvc is invoked\n", s)
	defer log.Printf("%T - InsertSocialSvc executed\n", s)

	if isValid, err := govalidator.ValidateStruct(social); !isValid {
		switch err.Error() {
		case "social  name is required":
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "SOCIAL_MEDIA_NAME_IS_EMPTY",
			}
			return result, errMsg
		case "social  url is required":
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "SOCIAL_MEDIA_URL_IS_EMPTY",
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

	social.UserID = userId

	log.Println("calling create social repo")
	result, err := s.socialRepo.InsertSocial(ctx, &social)

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

func (s *SocialUseCaseImpl) GetSocialsSvc(ctx context.Context) (result []social.Social, errMsg message.ErrorMessage) {
	log.Printf("%T - GetSocialsSvc is invoked\n", s)
	defer log.Printf("%T - GetSocialsSvc executed\n", s)

	stringUserId := ctx.Value("user").(string)

	userId, _ := strconv.ParseUint(stringUserId, 0, 64)

	log.Println("calling get social  by user id repo")
	result, err := s.socialRepo.GetSocials(ctx, userId)

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

func (s *SocialUseCaseImpl) GetSocialByIdSvc(ctx context.Context, socmedId uint64) (result social.Social, errMsg message.ErrorMessage) {
	log.Printf("%T - GetSocialByIdSvc is invoked\n", s)
	defer log.Printf("%T - GetSocialByIdSvc executed\n", s)

	log.Println("calling get social  by id repo")
	result, err := s.socialRepo.GetSocialById(ctx, socmedId)

	if result.ID <= 0 {
		log.Printf("social  with id %v not found", socmedId)

		err = fmt.Errorf("social  with id %v not found", socmedId)
		errMsg = message.ErrorMessage{
			Error: err,
			Type:  "SOCIAL_MEDIA_NOT_FOUND",
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

func (s *SocialUseCaseImpl) UpdateSocialSvc(ctx context.Context, inputSocial social.Social) (result social.Social, errMsg message.ErrorMessage) {
	log.Printf("%T - UpdateSocialSvc is invoked\n", s)
	defer log.Printf("%T - UpdateSocialSvc executed\n", s)

	if isValid, err := govalidator.ValidateStruct(inputSocial); !isValid {
		switch err.Error() {
		case "social  name is required":
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "SOCIAL_MEDIA_NAME_IS_EMPTY",
			}
			return result, errMsg
		case "social  url is required":
			errMsg := message.ErrorMessage{
				Error: err,
				Type:  "SOCIAL_MEDIA_URL_IS_EMPTY",
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

	result, err := s.socialRepo.UpdateSocial(ctx, inputSocial)

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

func (s *SocialUseCaseImpl) DeleteSocialSvc(ctx context.Context, socmedId uint64) (errMsg message.ErrorMessage) {
	log.Printf("%T - DeleteSocialSvc is invoked\n", s)
	defer log.Printf("%T - DeleteSocialSvc executed\n", s)

	log.Println("calling delete social repo")
	err := s.socialRepo.DeleteSocial(ctx, socmedId)

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

func NewSocialUseCase(socialRepo social.SocialRepo) social.SocialUseCase {
	return &SocialUseCaseImpl{socialRepo: socialRepo}
}
