package service

import (
	"fmt"
	"tours/model"
	"tours/repository"
)

type CheckpointCompletionService struct {
	CheckpointCompletionRepository *repository.CheckpointCompletionRepository
}

func (service *CheckpointCompletionService) Create(checkpointCompletion model.CheckpointCompletion) error {
	err := service.CheckpointCompletionRepository.Save(checkpointCompletion)
	if err != nil {
		return fmt.Errorf("failed to create checkpointCompletion: %w", err)
	}
	return nil
}

func (service *CheckpointCompletionService) GetByIDs(eid int, cid int) (*model.CheckpointCompletion, error) {
	checkpointCompletion, err := service.CheckpointCompletionRepository.FindByIds(eid, cid)
	if err != nil {
		return nil, fmt.Errorf("failed to get checkpointCompletion with IDs %d, %d: %w", eid, cid, err)
	}
	return checkpointCompletion, nil
}
