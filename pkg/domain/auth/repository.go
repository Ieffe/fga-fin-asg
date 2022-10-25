package auth

import (
	"context"
	"fin-asg/pkg/domain/user"
)

type AuthRepo interface {
	LoginUser(ctx context.Context, username string) (result user.User, err error)
}
