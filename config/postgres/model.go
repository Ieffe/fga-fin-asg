package postgres

import (
	"fin-asg/pkg/domain/comment"
	"fin-asg/pkg/domain/photo"
	"fin-asg/pkg/domain/social"
	"fmt"
	"log"
	"os/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// will be using GORM

// init config struct
type Config struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	DatabaseName string `json:"database_name"`
	User         string `json:"user"`
	Password     string `json:"password"`
}

// creating interface
type PostgresClient interface {
	GetClient() *gorm.DB
}

type PostgresClientImpl struct {
	cln    *gorm.DB
	config Config
}

func NewPostgresConnection(config Config) PostgresClient {
	connectionString := fmt.Sprintf(`
	host=%s 
	port=%s
	user=%s 
	password=%s 
	dbname=%s`,
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DatabaseName)
	
		
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		log.Fatal("error connection to database:", err)
	}

	fmt.Println("database connection is successfully connected")
	db.AutoMigrate(&user.User{}, &photo.Photo{}, &comment.Comment{}, &social.Social{})

	return &PostgresClientImpl{cln: db, config: config}
}

// implementation
func (p *PostgresClientImpl) GetClient() *gorm.DB {
	return p.cln
}
