package service

import (
	"fmt"
	"time"
	"tours/model"
	"tours/repository"
)

type TourExecutionService struct {
	TourExecutionRepository *repository.TourExecutionRepository
}

func (service *TourExecutionService) Create(uid int, tid int) (model.TourExecution, error) {
	//provera da li turu poseduje user
	tourExecution := model.TourExecution{
		TouristID:       uid,
		TourID:          tid,
		Start:           time.Now(),
		LastActivity:    time.Now(),
		ExecutionStatus: 2, //ili 0
	}

	err := service.TourExecutionRepository.Save(tourExecution)
	if err != nil {
		return tourExecution, fmt.Errorf("failed to create tourExecution: %w", err) //temp te
	}
	return tourExecution, nil //vratiti tur exec temp za sad
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
	err = service.TourExecutionRepository.Update(*tourExecution)
	if err != nil {
		return tourExecution, fmt.Errorf("failed to update tourExecution with ID %d: %w", tourExecution.ID, err)
	}
	return tourExecution, nil
}
