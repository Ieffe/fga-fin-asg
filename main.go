package main

import (
	"fin-asg/config/postgres"
	"fin-asg/pkg/domain/message"

	engine "fin-asg/config/gin"
	authrepo "fin-asg/pkg/repository/auth"
	authhandler "fin-asg/pkg/server/http/handler/auth"
	authusecase	"fin-asg/pkg/usecase/auth"
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
	userUsecase := userusecase.NewUserUsecase(userRepo)
	useHandler := userhandler.NewUserHandler(userUsecase)
	v1.NewUserRouter(ginEngine, useHandler).Routers()


	authRepo := authrepo.NewAuthRepo(postgresCln)
	authUsecase := authusecase.NewAuthUsecase(authRepo, userUsecase)
	authhandler := authhandler.NewAuthHandler(authUsecase)
	authMiddleware := middleware.NewAuthMiddleware(userUsecase)
	v1.NewAuthRouter(ginEngine, authhandler, authMiddleware).Routers()

	ginEngine.Serve()
}
