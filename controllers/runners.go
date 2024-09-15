package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runners/models"

	"github.com/gin-gonic/gin"
)

type Runners struct {
	runnersService *services.Runners
}

func NewRunners(runnersService *services.Runners) *Runners {
	return &Runners{
		runnersService: runnersService,
	}
}

func (r Runners) Create(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println("Error while unmarshalling create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response, responseErr := r.runnersService.Create(&runner)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (r Runners) Update(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading update runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println("Error while unmarshalling update runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	responseErr := r.runnersService.Update(&runner)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (r Runners) Delete(ctx *gin.Context) {
	runnerId := ctx.Param("id")
	responseErr := r.runnersService.Delete(runnerId)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (r Runners) Get(ctx *gin.Context) {
	runnerId := ctx.Param("id")
	response, responseErr := r.runnersService.Get(runnerId)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (r Runners) GetBatch(ctx *gin.Context) {
	params := ctx.Request.URL.Query()
	country := params.Get("country")
	year := params.Get("year")
	response, responseErr := r.runnersService.GetBatch(country, year)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
