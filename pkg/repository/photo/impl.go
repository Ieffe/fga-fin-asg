package photo

import (
	"context"
	"fin-asg/config/postgres"
	"fin-asg/pkg/domain/photo"
	"log"

	"gorm.io/gorm/clause"
)

type PhotoRepoImpl struct {
	pgCln postgres.PostgresClient
}

func (p *PhotoRepoImpl) InsertPhoto(ctx context.Context, input *photo.Photo) (result photo.Photo, err error) {
	log.Printf("%T - Inserthoto is invoked\n", p)
	defer log.Printf("%T - InsertPhoto executed\n", p)

	db := p.pgCln.GetClient()

	err = db.Model(&result).Create(&input).Error

	if err != nil {
		log.Printf("error when creating photo %v\n", input.ID)
	}

	result = *input

	return result, err
}

func (p *PhotoRepoImpl) GetPhotosByUserId(ctx context.Context, userId uint64) (result []photo.Photo, err error) {
	log.Printf("%T - GetPhotosByUserId is invoked\n", p)
	defer log.Printf("%T - GetPhotosByUserId executed\n", p)

	db := p.pgCln.GetClient()

	err = db.Table("photos").Where("user_id = ?", userId).Select("id", "created_at", "updated_at", "title", "caption", "url", "user_id").Order("id").Preload("Comments").Find(&result).Error

	if err != nil {
		log.Printf("error when getting photos by user id %v\n", userId)
	}

	return result, err
}

func (p *PhotoRepoImpl) GetPhotoById(ctx context.Context, photoId uint64) (result photo.Photo, err error) {
	log.Printf("%T - GetPhotoById is invoked\n", p)
	defer log.Printf("%T - GetPhotoById executed\n", p)

	db := p.pgCln.GetClient()

	err = db.Table("photos").Where("id = ?", photoId).Select("id", "title", "caption", "url", "user_id").Preload("Comments").Find(&result).Error

	if err != nil {
		log.Printf("error when getting photo by id %v\n", photoId)
	}

	return result, err
}

func (p *PhotoRepoImpl) UpdatePhoto(ctx context.Context, title string, caption string, url string) (result photo.Photo, err error) {
	log.Printf("%T - UpdatePhoto is invoked\n", p)
	defer log.Printf("%T - UpdatePhoto executed\n", p)

	photoId := ctx.Value("photoId").(uint64)

	db := p.pgCln.GetClient()

	err = db.Model(&result).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}, {Name: "title"}, {Name: "caption"}, {Name: "url"}, {Name: "user_id"}, {Name: "updated_at"}}}).Where("id = ?", photoId).Updates(photo.Photo{Title: title, Caption: caption, Url: url}).Error

	if err != nil {
		log.Printf("error when updating photo by id %v\n", photoId)
	}

	return result, err
}

func (p *PhotoRepoImpl) DeletePhoto(ctx context.Context, photoId uint64) (err error) {
	log.Printf("%T - DeletePhoto is invoked\n", p)
	defer log.Printf("%T - DeletePhoto executed\n", p)

	db := p.pgCln.GetClient()

	err = db.Where("id = ?", photoId).Delete(&photo.Photo{}).Error

	if err != nil {
		log.Printf("error when deleting photo by id %v \n", photoId)
	}

	return err
}

func NewPhotoRepo(pgCln postgres.PostgresClient) photo.PhotoRepo {
	return &PhotoRepoImpl{pgCln: pgCln}
}
