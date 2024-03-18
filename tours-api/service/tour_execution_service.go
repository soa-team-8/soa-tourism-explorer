package service

import (
	"fmt"
	"math"
	"time"
	"tours/model"
	"tours/repository"
)

type TourExecutionService struct {
	TourExecutionRepository *repository.TourExecutionRepository
}

func (service *TourExecutionService) Create(userID uint64, tourId uint64) (model.TourExecution, error) {
	//TODO: check if user owns the tour
	tourExecution := model.TourExecution{
		TouristID:       userID,
		TourID:          tourId,
		Start:           time.Now(),
		LastActivity:    time.Now(),
		ExecutionStatus: model.InProgress,
	}

	err := service.TourExecutionRepository.Save(tourExecution)
	if err != nil {
		return model.TourExecution{}, fmt.Errorf("failed to create tourExecution: %w", err)
	}
	newTourExecution, err := service.GetByIDs(userID, tourId)
	if err != nil {
		return model.TourExecution{}, fmt.Errorf("failed to get tourExecution with IDs %d, %d: %w", userID, tourId, err)
	}
	return *newTourExecution, nil
}

func (service *TourExecutionService) GetByIDs(userID uint64, tourID uint64) (*model.TourExecution, error) {
	executionExists, err := service.TourExecutionRepository.ExistsByIDs(userID, tourID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tourExecution with IDs %d, %d: %w", userID, tourID, err)
	}
	if !executionExists {
		newTourExecution, err := service.Create(userID, tourID)
		if err != nil {
			return nil, fmt.Errorf("failed to create tourExecution: %w", err)
		}
		return &newTourExecution, nil
	}
	tourExecution, err := service.TourExecutionRepository.FindInProgressByIds(userID, tourID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tourExecution with IDs %d, %d: %w", userID, tourID, err)
	}
	return tourExecution, nil
}

func (service *TourExecutionService) Abandon(userID uint64, executionID uint64) (*model.TourExecution, error) {
	tourExecution, err := service.TourExecutionRepository.FindByID(executionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tourExecution with ID %d: %w", executionID, err)
	}
	if tourExecution == nil {
		return nil, fmt.Errorf("tourExecution with ID %d not found", executionID)
	}
	//TODO: check if user owns the tour
	if tourExecution.ExecutionStatus != model.InProgress {
		return &model.TourExecution{}, fmt.Errorf("can not abandon tourExecution with ID %d: %w", executionID, err)
	}
	tourExecution.ExecutionStatus = model.Abandoned
	tourExecution.LastActivity = time.Now()
	err = service.TourExecutionRepository.Update(*tourExecution)
	if err != nil {
		return &model.TourExecution{}, fmt.Errorf("failed to update tourExecution with ID %d: %w", tourExecution.ID, err)
	}
	return tourExecution, nil
}

func (service *TourExecutionService) CheckPosition(touristPosition model.TouristPosition, id uint64) (*model.TourExecution, error) {
	tourExecution, err := service.TourExecutionRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tourExecution with ID %d: %w", id, err)
	}

	if tourExecution.ExecutionStatus == model.Abandoned {
		return tourExecution, nil
	}

	index := 0
	for _, checkpoint := range tourExecution.Tour.Checkpoints {
		a := math.Abs(math.Round(checkpoint.Longitude*10000)/10000 - math.Round(touristPosition.Longitude*10000)/10000)
		b := math.Abs(math.Round(checkpoint.Latitude*10000)/10000 - math.Round(touristPosition.Latitude*10000)/10000)

		if a < 0.004 && b < 0.004 {
			executionID := tourExecution.ID
			checkpointID := tourExecution.Tour.Checkpoints[index].ID
			completion := model.CheckpointCompletion{
				TourExecutionID: executionID,
				CheckpointID:    checkpointID,
				CompletionTime:  time.Now(),
			}

			completionExists := false
			for _, cp := range tourExecution.CompletedCheckpoints {
				if cp.TourExecutionID == completion.TourExecutionID && cp.CheckpointID == completion.CheckpointID {
					completionExists = true
					break
				}
			}

			if !completionExists {
				tourExecution.CompletedCheckpoints = append(tourExecution.CompletedCheckpoints, completion)
				err = service.TourExecutionRepository.Update(*tourExecution)
				if err != nil {
					return tourExecution, fmt.Errorf("failed to update tourExecution with ID %d: %w", tourExecution.ID, err)
				}
			}
		}
		index++
	}

	if len(tourExecution.CompletedCheckpoints) == len(tourExecution.Tour.Checkpoints) {
		tourExecution.ExecutionStatus = model.Completed
	}
	tourExecution.LastActivity = time.Now()
	err = service.TourExecutionRepository.Update(*tourExecution)
	if err != nil {
		return tourExecution, fmt.Errorf("failed to update tourExecution with ID %d: %w", tourExecution.ID, err)
	}
	return tourExecution, nil
}
