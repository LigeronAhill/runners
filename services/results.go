package services

import (
	"log"
	"net/http"
	"runners/models"
	"runners/repositories"
	"time"
)

type Results struct {
	resultsRepository *repositories.Results
	runnersRepository *repositories.Runners
}

func NewResults(
	resultsRepository *repositories.Results,
	runnersRepository *repositories.Runners,
) *Results {
	return &Results{
		resultsRepository: resultsRepository,
		runnersRepository: runnersRepository,
	}
}

func (r Results) Create(result *models.Result) (*models.Result, *models.ResponseError) {
	if result.RunnerID == "" {
		return nil, &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}
	if result.RaceResult == "" {
		return nil, &models.ResponseError{
			Message: "Invalid race result",
			Status:  http.StatusBadRequest,
		}
	}
	if result.Location == "" {
		return nil, &models.ResponseError{
			Message: "Invalid location",
			Status:  http.StatusBadRequest,
		}
	}
	if result.Position < 0 {
		return nil, &models.ResponseError{
			Message: "Invalid position",
			Status:  http.StatusBadRequest,
		}
	}
	currentYear := time.Now().Year()
	if result.Year < 0 || result.Year > currentYear {
		return nil, &models.ResponseError{
			Message: "Invalid year",
			Status:  http.StatusBadRequest,
		}
	}
	raceResult, err := parseRaceResult(result.RaceResult)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid race result",
			Status:  http.StatusBadRequest,
		}
	}
	response, responseErr := r.resultsRepository.Create(result)
	if responseErr != nil {
		return nil, responseErr
	}
	runner, responseErr := r.runnersRepository.Get(result.RunnerID)
	if responseErr != nil {
		return nil, responseErr
	}
	if runner == nil {
		return nil, &models.ResponseError{
			Message: "Runner not found",
			Status:  http.StatusNotFound,
		}
	}
	if runner.PersonalBest == "" {
		runner.PersonalBest = result.RaceResult
	} else {
		personalBest, err := parseRaceResult(runner.PersonalBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "Failed to parse personal best",
				Status:  http.StatusInternalServerError,
			}
		}
		if raceResult < personalBest {
			runner.PersonalBest = result.RaceResult
		}
	}
	if result.Year == currentYear {
		if runner.SeasonBest == "" {
			runner.SeasonBest = result.RaceResult
		} else {
			seasonBest, err := parseRaceResult(runner.SeasonBest)
			if err != nil {
				return nil, &models.ResponseError{
					Message: "Failed to parse season best",
					Status:  http.StatusInternalServerError,
				}
			}
			if raceResult < seasonBest {
				runner.SeasonBest = result.RaceResult
			}
		}
	}
	responseErr = r.runnersRepository.UpdateResults(runner)
	if responseErr != nil {
		return nil, responseErr
	}
	return response, nil
}

func (r Results) Delete(resultId string) *models.ResponseError {
	if resultId == "" {
		return &models.ResponseError{
			Message: "Invalid result ID",
			Status:  http.StatusBadRequest,
		}
	}
	result, responseErr := r.resultsRepository.Delete(resultId)
	if responseErr != nil {
		return responseErr
	}
	runner, responseErr := r.runnersRepository.Get(result.RunnerID)
	if responseErr != nil {
		return responseErr
	}
	if runner.PersonalBest == result.RaceResult {
		personalBest, responseErr := r.resultsRepository.GetPersonalBestResults(result.RunnerID)
		if responseErr != nil {
			return responseErr
		}
		runner.PersonalBest = personalBest
	}
	currentYear := time.Now().Year()
	if runner.SeasonBest == result.RaceResult && result.Year == currentYear {
		seasonBest, responseErr := r.resultsRepository.GetSeasonBestResults(result.RunnerID, result.Year)
		if responseErr != nil {
			return responseErr
		}
		runner.SeasonBest = seasonBest
	}
	responseErr = r.runnersRepository.UpdateRunnerResults(runner)
	if responseErr != nil {
		return responseErr
	}
	return nil
}

func parseRaceResult(timeString string) (time.Duration, error) {
	return time.ParseDuration(
		timeString[0:2] + "h" +
			timeString[3:5] + "m" +
			timeString[6:8] + "s",
	)
}
