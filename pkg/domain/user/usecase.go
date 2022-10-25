package user

import (
	"context"
	"fin-asg/pkg/domain/message"
)

type UserUseCase interface {
	GetUserByIdSvc(ctx context.Context, id uint64) (result User, errMsg message.ErrorMessage)
	InsertUserSvc(ctx context.Context, input User) (result User, errMsg message.ErrorMessage)
	UpdateUserSvc(ctx context.Context, id uint64, email string, username string) (idToken string, errMsg message.ErrorMessage)
	DeleteUserSvc(ctx context.Context, id uint64) (err message.ErrorMessage)
}