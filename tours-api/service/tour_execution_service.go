package service

import (
	"fmt"
	"math"
	"time"
	"tours/model"
	"tours/repository"
)

type TourExecutionService struct {
	TourExecutionRepository        *repository.TourExecutionRepository
	CheckpointCompletionRepository *repository.CheckpointCompletionRepository
}

func (service *TourExecutionService) Create(uid int, tid int) (model.TourExecution, error) {
	//provera da li turu poseduje user
	tourExecution := model.TourExecution{
		TouristID: uint64(uid),
		//TourID:          tid,
		Start:           time.Now(),
		LastActivity:    time.Now(),
		ExecutionStatus: 2,
	}

	err := service.TourExecutionRepository.Save(tourExecution)
	if err != nil {
		return tourExecution, fmt.Errorf("failed to create tourExecution: %w", err) //temp te
	}

	newTourExecution, err := service.GetByIDs(uid, tid)
	if err != nil {
		return model.TourExecution{}, fmt.Errorf("failed to get tourExecution with IDs %d, %d: %w", uid, tid, err)
	}
	return *newTourExecution, nil //testirati
}

func (service *TourExecutionService) GetByIDs(uid int, tid int) (*model.TourExecution, error) {
	tourExecution, err := service.TourExecutionRepository.FindInProgressByIds(uid, tid)
	if err != nil {
		return nil, fmt.Errorf("failed to get tourExecution with IDs %d, %d: %w", uid, tid, err)
	}
	return tourExecution, nil
}

func (service *TourExecutionService) Abandon(uid int, eid int) (*model.TourExecution, error) {
	tourExecution, err := service.TourExecutionRepository.FindByID(eid)
	if err != nil {
		return nil, fmt.Errorf("failed to get tourExecution with ID %d: %w", eid, err)
	}
	if tourExecution == nil {
		return nil, fmt.Errorf("tourExecution with ID %d not found", eid)
	}
	//provera dal je kupio
	tourExecution.ExecutionStatus = model.Abandoned
	tourExecution.LastActivity = time.Now()
	err = service.TourExecutionRepository.Update(*tourExecution)
	if err != nil {
		return tourExecution, fmt.Errorf("failed to update tourExecution with ID %d: %w", tourExecution.ID, err)
	}
	return tourExecution, nil
}

func (service *TourExecutionService) CheckPosition(touristPosition model.TouristPosition, id int) (*model.TourExecution, error) {
	tourExecution, err := service.TourExecutionRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tourExecution with ID %d: %w", id, err)
	} //dodati da povuce sve checkpointe u okviru ture

	checkpoints := tourExecution.Tour.Checkpoints
	indexNumber := 0
	for _, checkpoint := range checkpoints {
		a := math.Abs(math.Round(checkpoint.Longitude*10000)/10000 - math.Round(touristPosition.Longitude*10000)/10000)
		b := math.Abs(math.Round(checkpoint.Latitude*10000)/10000 - math.Round(touristPosition.Latitude*10000)/10000)

		if a < 0.01 && b < 0.01 {
			eid := int64(tourExecution.ID)
			cid := int64(tourExecution.Tour.Checkpoints[indexNumber].ID)
			completion := model.CheckpointCompletion{
				TourExecutionID: uint64(eid),
				CheckpointID:    uint64(cid),
				CompletionTime:  time.Now(),
			}
			_, err := service.CheckpointCompletionRepository.FindByIds(int(eid), int(cid))
			if err == nil {
				err := service.CheckpointCompletionRepository.Save(completion)
				if err != nil {
					return nil, fmt.Errorf("failed to save completion with ID %d: %w", id, err)
				}
			}
		}
		indexNumber++
	}

	tourExecution, err = service.TourExecutionRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tourExecution with ID %d: %w", id, err)
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
