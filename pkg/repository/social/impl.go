package social

import (
	"context"
	"fin-asg/config/postgres"
	"fin-asg/pkg/domain/social"
	"log"

	"gorm.io/gorm/clause"
)

type SocialRepoImpl struct {
	pgCln postgres.PostgresClient
}

func (s *SocialRepoImpl) InsertSocial(ctx context.Context, inputSocial *social.Social) (result social.Social, err error) {
	log.Printf("%T - InsertSocial is invoked\n", s)
	defer log.Printf("%T - InsertSocial executed\n", s)

	db := s.pgCln.GetClient()

	err = db.Model(&result).Create(&inputSocial).Error

	if err != nil {
		log.Printf("error when creating social for photo id %v\n", inputSocial)
	}

	result = *inputSocial

	return result, err
}

func (s *SocialRepoImpl) GetSocials(ctx context.Context, userId uint64) (result []social.Social, err error) {
	log.Printf("%T - GetSocials is invoked\n", s)
	defer log.Printf("%T - GetSocials executed\n", s)

	db := s.pgCln.GetClient()

	err = db.Model(&social.Social{}).Where("user_id = ?", userId).Find(&result).Error

	if err != nil {
		log.Printf("error when getting social  by user id %v\n", userId)
	}

	return result, err
}

func (s *SocialRepoImpl) GetSocialById(ctx context.Context, socmedId uint64) (result social.Social, err error) {
	log.Printf("%T - GetSocialById is invoked\n", s)
	defer log.Printf("%T - GetSocialById executed\n", s)

	db := s.pgCln.GetClient()

	err = db.Table("social_").Where("id = ?", socmedId).Find(&result).Error

	if err != nil {
		log.Printf("error when getting social  by id %v\n", socmedId)
	}

	return result, err
}

func (s *SocialRepoImpl) UpdateSocial(ctx context.Context, inputSocial social.Social) (result social.Social, err error) {
	log.Printf("%T - UpdateSocial is invoked\n", s)
	defer log.Printf("%T - UpdateSocial executed\n", s)

	id:= ctx.Value("id").(uint64)

	db := s.pgCln.GetClient()

	err = db.Model(&result).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}, {Name: "name"}, {Name: "url"}, {Name: "user_id"}, {Name: "updated_at"}}}).Where("id = ?", id).Updates(inputSocial).Error

	if err != nil {
		log.Printf("error when updating social by id %v\n", id)
	}

	return result, err
}

func (s *SocialRepoImpl) DeleteSocial(ctx context.Context, socmedId uint64) (err error) {
	log.Printf("%T - DeleteSocial is invoked\n", s)
	defer log.Printf("%T - DeleteSocial executed\n", s)

	db := s.pgCln.GetClient()

	err = db.Where("id = ?", socmedId).Delete(&social.Social{}).Error

	if err != nil {
		log.Printf("error when deleting social  by id %v \n", socmedId)
	}

	return err
}

func NewSocialRepo(pgCln postgres.PostgresClient) social.SocialRepo {
	return &SocialRepoImpl{pgCln: pgCln}
}
