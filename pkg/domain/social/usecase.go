package social

import (
	"context"
	"fin-asg/pkg/domain/message"
)

type SocialUseCase interface {
	InsertSocialSvc(ctx context.Context, social Social) (result Social, errMsg message.ErrorMessage)
	GetSocialsSvc(ctx context.Context) (result []Social, errMsg message.ErrorMessage)
	GetSocialByIdSvc(ctx context.Context, socmedId uint64) (result Social, errMsg message.ErrorMessage)
	UpdateSocialSvc(ctx context.Context, social Social) (result Social, errMsg message.ErrorMessage)
	DeleteSocialSvc(ctx context.Context, socmedId uint64) (errMsg message.ErrorMessage)
}