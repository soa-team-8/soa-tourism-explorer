package service

import (
	"fmt"
	"time"
	"tours/model"
	"tours/repository"
)

type TourRatingService struct {
	TourRatingRepository    *repository.TourRatingRepository
	TourExecutionRepository *repository.TourExecutionRepository
}

func (service *TourRatingService) validateTourRating(tourRating model.TourRating) error {
	if tourRating.Rating < 1 || tourRating.Rating > 5 {
		return fmt.Errorf("failed to update: rating can not be %d", tourRating.Rating)
	}

	tourExecution, err := service.TourExecutionRepository.FindByIds(tourRating.TouristID, tourRating.TourID)
	if err != nil {
		return fmt.Errorf("failed to get tourExecution: %w", err)
	}

	totalCheckpoints := len(tourExecution.Tour.Checkpoints)
	completedCheckpoints := len(tourExecution.CompletedCheckpoints)
	completionPercentage := float64(completedCheckpoints) / float64(totalCheckpoints) * 100
	if completionPercentage < 35 {
		return fmt.Errorf("you have not completed 35 perc of the tour (%f)", completionPercentage)
	}

	lastActivityTime := tourExecution.LastActivity
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	if lastActivityTime.Before(oneWeekAgo) {
		return fmt.Errorf("more than a week has passed since the tour was activated")
	}

	return nil
}

func (service *TourRatingService) Create(tourRating model.TourRating) (uint64, error) {
	isRated, err := service.TourRatingRepository.ExistsByIDs(tourRating.TouristID, tourRating.TourID)
	if err != nil {
		return 0, nil
	}
	if isRated {
		return 0, fmt.Errorf("you already rated this tour")
	}

	if err = service.validateTourRating(tourRating); err != nil {
		return 0, err
	}

	id, err := service.TourRatingRepository.Save(tourRating)
	if err != nil {
		return 0, fmt.Errorf("failed to create tourRating: %w", err)
	}
	return id, nil
}

func (service *TourRatingService) Update(tourRating model.TourRating) error {
	if err := service.validateTourRating(tourRating); err != nil {
		return err
	}

	if err := service.TourRatingRepository.Update(tourRating); err != nil {
		return fmt.Errorf("failed to update tourRating with ID %d: %w", tourRating.ID, err)
	}
	return nil
}

func (service *TourRatingService) Delete(id uint64) error {
	if err := service.TourRatingRepository.Delete(id); err != nil {
		return fmt.Errorf("failed to delete tourRating with ID %d: %w", id, err)
	}
	return nil
}

func (service *TourRatingService) GetAll() ([]model.TourRating, error) {
	tourRating, err := service.TourRatingRepository.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all tourRatings: %w", err)
	}
	return tourRating, nil
}

func (service *TourRatingService) GetByID(id uint64) (*model.TourRating, error) {
	tourRating, err := service.TourRatingRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tourRating with ID %d: %w", id, err)
	}
	return tourRating, nil
}
