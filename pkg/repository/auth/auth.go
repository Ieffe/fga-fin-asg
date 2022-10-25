package auth

import (
	"context"
	"fin-asg/config/postgres"
	"fin-asg/pkg/domain/auth"
	"fin-asg/pkg/domain/user"
	"log"
)

type AuthRepoImpl struct {
	pgCln postgres.PostgresClient
}

func (a *AuthRepoImpl) LoginUser(ctx context.Context, username string) (result user.User, err error) {
	log.Printf("%T - LoginUser is invoked\n", a)
	defer log.Printf("%T - LoginUser executed\n", a)

	db := a.pgCln.GetClient()

	err = db.Model(&user.User{}).Select("id", "username", "password", "email", "dob").Where("username = ?", username).Find(&result).Error

	if err != nil {
		log.Printf("error when getting username %v\n", username)
	}

	return result, err
}

func NewAuthRepo(pgCln postgres.PostgresClient) auth.AuthRepo {
	return &AuthRepoImpl{pgCln: pgCln}
}