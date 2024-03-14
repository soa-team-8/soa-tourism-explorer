package service

import (
	"errors"
	"fmt"
	"tours/model"
	"tours/repository"
)

type TourRatingService struct {
	TourRatingRepository *repository.TourRatingRepository
}

func (service *TourRatingService) Create(tourRating model.TourRating) error {
	if tourRating.Rating < 1 || tourRating.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	err := service.TourRatingRepository.Save(tourRating)
	if err != nil {
		return fmt.Errorf("failed to create tourRating: %w", err)
	}
	return nil
}

func (service *TourRatingService) Delete(id uint64) error {
	if err := service.TourRatingRepository.Delete(id); err != nil {
		return fmt.Errorf("failed to delete tourRating with ID %d: %w", id, err)
	}
	return nil
}

// krijejt imidz aploud i apdejt tur egzekjusn validacija
func (service *TourRatingService) Update(tourRating model.TourRating) error {
	if tourRating.Rating < 1 || tourRating.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	if err := service.TourRatingRepository.Update(tourRating); err != nil {
		return fmt.Errorf("failed to update tourRating with ID %d: %w", tourRating.ID, err)
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
