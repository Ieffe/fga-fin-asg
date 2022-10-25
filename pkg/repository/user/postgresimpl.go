package user

import (
	"context"
	"fin-asg/config/postgres"
	"fin-asg/pkg/domain/user"
	"log"

	"gorm.io/gorm/clause"
)

type UserRepoImpl struct {
	pgCln postgres.PostgresClient
}

func NewUserRepo(pgCln postgres.PostgresClient) user.UserRepo {
	return &UserRepoImpl{pgCln: pgCln}
}

func (u *UserRepoImpl) GetUserById(ctx context.Context, id uint64) (result user.User, err error) {
	log.Printf("%T - GetUserById is invoked\n", u)
	defer log.Printf("%T - GetUserById executed\n", u)

	db := u.pgCln.GetClient()

	err = db.Model(&user.User{}).Where("id = ?", id).Preload("SocialMedias").Find(&result).Error

	if err != nil {
		log.Printf("error when getting user by id %v\n", id)
	}

	return result, err
}

func (u *UserRepoImpl) InsertUser(ctx context.Context, insertedUser *user.User) (err error) {
	log.Printf("%T - InsertUser is invoked\n", u)
	defer log.Printf("%T - InsertUser executed\n", u)

	db := u.pgCln.GetClient()

	err = db.Model(&user.User{}).Create(&insertedUser).Error

	if err != nil {
		log.Printf("error when creating user %v\n", insertedUser.Email)
	}

	return err
}

func (u *UserRepoImpl) UpdateUser(ctx context.Context, id uint64, email string, username string) (result user.User, err error) {
	log.Printf("%T - UpdateUser is invoked\n", u)
	defer log.Printf("%T - UpdateUser executed\n", u)

	db := u.pgCln.GetClient()

	err = db.Model(&result).Clauses(clause.Returning{}).Where("id = ?", id).Updates(user.User{Email: email, Username: username}).Error

	if err != nil {
		log.Printf("error when updating user by id %v\n", id)
	}

	return result, err
}

func (u *UserRepoImpl) DeleteUser(ctx context.Context, id uint64) (err error) {
	log.Printf("%T - DeleteUser is invoked\n", u)
	defer log.Printf("%T - DeleteUser executed\n", u)

	db := u.pgCln.GetClient()

	err = db.Where("id = ?",  id).Delete(&user.User{}).Error

	if err != nil {
		log.Printf("error when deleting user by id %v \n", id)
	}

	return err
}
