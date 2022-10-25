package auth

import (
	"context"
	"fin-asg/pkg/domain/message"
)

type AuthUsecase interface {
	LoginUserSvc(ctx context.Context, username string, password string) (accessToken string, refreshToken string, idToken string, errMsg message.ErrorMessage)
	GetRefreshTokenSvc(ctx context.Context) (accessToken string, refreshToken string, idToken string, errMsg message.ErrorMessage)
}
