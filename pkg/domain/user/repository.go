package user

import "context"

type UserRepo interface {
	GetUserById(ctx context.Context, id uint64) (result User, err error)
	InsertUser(ctx context.Context, insertedUser *User) (err error)
	UpdateUser(ctx context.Context, id uint64, email string, username string) (result User, err error)
	DeleteUser(ctx context.Context, id uint64) (err error)
}