package social

import "context"

type SocialRepo interface {
	InsertSocial(ctx context.Context, social *Social) (result Social, err error)
	GetSocials(ctx context.Context, userId uint64) (result []Social, err error)
	GetSocialById(ctx context.Context, id uint64) (result Social, err error)
	UpdateSocial(ctx context.Context, social Social) (result Social, err error)
	DeleteSocial(ctx context.Context, socmedId uint64) (err error)
}