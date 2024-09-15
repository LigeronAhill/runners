package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runners/models"
	"runners/services"

	"github.com/gin-gonic/gin"
)

type Results struct {
	resultsService *services.Results
}

func NewResults(resultsService *services.Results) *Results {
	return &Results{resultsService}
}

func (r Results) Create(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var result models.Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("error while unmarshalling vreate result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response, responseErr := r.resultsService.Create(&result)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (r Results) Delete(ctx *gin.Context) {
	resultId := ctx.Param("id")
	responseErr := r.resultsService.Delete(resultId)
	if responseErr != nil {
		ctx.JSON(responseErr.Status, responseErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}
