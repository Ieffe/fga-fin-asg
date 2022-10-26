package main

import (
	"fin-asg/config/postgres"
	"fin-asg/pkg/domain/message"

	engine "fin-asg/config/gin"

	authrepo "fin-asg/pkg/repository/auth"
	authhandler "fin-asg/pkg/server/http/handler/auth"
	authusecase	"fin-asg/pkg/usecase/auth"

	photorepo  "fin-asg/pkg/repository/photo"
	photohandler "fin-asg/pkg/server/http/handler/photo"
	photousecase "fin-asg/pkg/usecase/photo"

	commentrepo  "fin-asg/pkg/repository/comment"
	commenthandler "fin-asg/pkg/server/http/handler/comment"
	commentusecase "fin-asg/pkg/usecase/comment"

	socialrepo  "fin-asg/pkg/repository/social"
	socialhandler "fin-asg/pkg/server/http/handler/social"
	socialusecase "fin-asg/pkg/usecase/social"


	userrepo "fin-asg/pkg/repository/user"
	userhandler "fin-asg/pkg/server/http/handler/user"
	userusecase "fin-asg/pkg/usecase/user"
	"fin-asg/pkg/server/http/middleware"
	v1 "fin-asg/pkg/server/http/router/v1"
	

	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	postgresCln := postgres.NewPostgresConnection(postgres.Config{
		Host:         goDotEnvVariable("HOST"),
		Port:         goDotEnvVariable("PORT"),
		User:         goDotEnvVariable("USER"),
		Password:     goDotEnvVariable("PASSWORD"),
		DatabaseName: goDotEnvVariable("DB_NAME"),
	})

	ginEngine := engine.NewGinHttp(engine.Config{
		Port: ":8080",
	})

	ginEngine.GetGin().Use(
		gin.Recovery(),
		gin.Logger())

	startTime := time.Now()
	ginEngine.GetGin().GET("/", func(ctx *gin.Context) {
		// secara default map jika di return dalam
		// response API, dia akan menjadi JSON
		respMap := map[string]any{
			"code":       0,
			"message":    "server up and running",
			"start_time": startTime,
		}

		// golang memiliki json package
		// json package bisa mentranslasikan
		// map menjadi suatu struct
		// nb: struct harus memiliki tag/annotation JSON
		var respStruct message.Response

		// marshal -> mengubah json/struct/map menjadi
		// array of byte atau bisa kita translatekan menjadi string
		// dengan format JSON
		resByte, err := json.Marshal(respMap)
		if err != nil {
			panic(err)
		}
		// unmarshal -> translasikan string/[]byte dengan format JSON
		// menjadi map/struct dengan tag/annotation json
		err = json.Unmarshal(resByte, &respStruct)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, respStruct)
	})
	
	userRepo := userrepo.NewUserRepo(postgresCln)
	userUseCase := userusecase.NewUserUsecase(userRepo)
	useHandler := userhandler.NewUserHandler(userUseCase)
	v1.NewUserRouter(ginEngine, useHandler).Routers()

	authRepo := authrepo.NewAuthRepo(postgresCln)
	authUseCase := authusecase.NewAuthUsecase(authRepo, userUseCase)
	authhandler := authhandler.NewAuthHandler(authUseCase)
	authMiddleware := middleware.NewAuthMiddleware(userUseCase)

	photoRepo := photorepo.NewPhotoRepo(postgresCln)
	photoUseCase := photousecase.NewPhotoUsecase(photoRepo, userUseCase)
	photoHandler := photohandler.NewPhotoHandler(photoUseCase)
	v1.NewPhotoRouter(ginEngine, photoHandler, authMiddleware).Routers()

	commentRepo := commentrepo.NewCommentRepo(postgresCln)
	commentUseCase := commentusecase.NewCommentUsecase(commentRepo, photoUseCase)
	commentHandler := commenthandler.NewCommentHandler(commentUseCase)
	v1.NewCommentRouter(ginEngine, commentHandler, authMiddleware).Routers()

	socialRepo := socialrepo.NewSocialRepo(postgresCln)
	socialUseCase := socialusecase.NewSocialUseCase(socialRepo)
	socialHandler := socialhandler.NewSocialHandler(socialUseCase)
	v1.NewSocialRouter(ginEngine, socialHandler, authMiddleware).Routers()
	
	v1.NewAuthRouter(ginEngine, authhandler, authMiddleware).Routers()

	ginEngine.Serve()
}
