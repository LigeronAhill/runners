package services

import (
	"net/http"
	"runners/models"
	"runners/repositories"
	"strconv"
	"time"
)

type Runners struct {
	runnersRepository *repositories.Runners
	resultsRepository *repositories.Results
}

func NewRunners(
	runnersRepository *repositories.Runners,
	resultsRepository *repositories.Results,
) *Runners {
	return &Runners{
		runnersRepository: runnersRepository,
		resultsRepository: resultsRepository,
	}
}

func (r Runners) Create(runner *models.Runner) (*models.Runner, *models.ResponseError) {
	responseErr := validate(runner)
	if responseErr != nil {
		return nil, responseErr
	}
	return r.runnersRepository.Create(runner)
}

func (r Runners) Update(runner *models.Runner) *models.ResponseError {
	responseErr := validateId(runner.ID)
	if responseErr != nil {
		return responseErr
	}
	responseErr = validate(runner)
	if responseErr != nil {
		return responseErr
	}
	return r.runnersRepository.Update(runner)
}

func (r Runners) Delete(runnerId string) *models.ResponseError {
	responseErr := validateId(runnerId)
	if responseErr != nil {
		return responseErr
	}
	return r.runnersRepository.Delete(runnerId)
}

func (r Runners) Get(runnerId string) (*models.Runner, *models.ResponseError) {
	responseErr := validateId(runnerId)
	if responseErr != nil {
		return nil, responseErr
	}
	runner, responseErr := r.runnersRepository.Get(runnerId)
	if responseErr != nil {
		return nil, responseErr
	}
	results, responseErr := r.resultsRepository.GetAllRunnersResults(runnerId)
	if responseErr != nil {
		return nil, responseErr
	}
	runner.Results = results
	return runner, nil
}

func (r Runners) GetBatch(country string, year string) ([]*models.Runner, *models.ResponseError) {
	if country != "" && year != "" {
		return nil, &models.ResponseError{
			Message: "Only one parameter can be passed",
			Status:  http.StatusBadRequest,
		}
	}
	if country != "" {
		return r.runnersRepository.GetRunnerByCountry(country)
	}
	if year != "" {
		intYear, err := strconv.Atoi(year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "Invalid year",
				Status:  http.StatusBadRequest,
			}
		}
		currentYear := time.Now().Year()
		if intYear < 0 || intYear > currentYear {
			return nil, &models.ResponseError{
				Message: "Invalid year",
				Status:  http.StatusBadRequest,
			}
		}
		return r.runnersRepository.GetRunnersByYear(intYear)
	}
	return r.runnersRepository.GetAllRunners()
}

func validate(runner *models.Runner) *models.ResponseError {
	if runner.FirstName == "" {
		return &models.ResponseError{
			Message: "Invalid first name",
			Status:  http.StatusBadRequest,
		}
	}
	if runner.LastName == "" {
		return &models.ResponseError{
			Message: "Invalid last name",
			Status:  http.StatusBadRequest,
		}
	}
	if runner.Age <= 16 || runner.Age > 125 {
		return &models.ResponseError{
			Message: "Invalid age",
			Status:  http.StatusBadRequest,
		}
	}
	if runner.Country == "" {
		return &models.ResponseError{
			Message: "Invalid country",
			Status:  http.StatusBadRequest,
		}
	}
	return nil
}

func validateId(runnerId string) *models.ResponseError {
	if runnerId == "" {
		return &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}
	return nil
}
