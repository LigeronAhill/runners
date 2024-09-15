package server

import (
	"database/sql"
	"log"
	"runners/controllers"
	"runners/repositories"
	"runners/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpServer struct {
	config            *viper.Viper
	router            *gin.Engine
	runnersController *controllers.Runners
	resultsController *controllers.Results
}

func InitHttp(config *viper.Viper, dbHandler *sql.DB) HttpServer {
	runnersRepository := repositories.NewRunners(dbHandler)
	resultsRepository := repositories.NewResults(dbHandler)
	runnersService := services.NewRunners(runnersRepository, resultsRepository)
	resultsService := services.NewResults(resultsRepository, runnersRepository)
	runnersController := controllers.NewRunners(runnersService)
	resultsController := controllers.NewResults(resultsService)
	router := gin.Default()
	router.POST("/runner", runnersController.Create)
	router.PUT("/runner", runnersController.Update)
	router.DELETE("/runner/:id", runnersController.Delete)
	router.GET("/runner/:id", runnersController.Get)
	router.GET("/runner", runnersController.GetBatch)
	router.POST("/result", resultsController.Create)
	router.DELETE("result/:id", resultsController.Delete)
	return HttpServer{
		config,
		router,
		runnersController,
		resultsController,
	}
}

func (hs HttpServer) Start() {
	err := hs.router.Run(hs.config.GetString("http.server_address"))
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
