package service

import (
	"fmt"
	"tours/model"
	"tours/repository"
)

type TourService struct {
	TourRepository *repository.TourRepository
}

func (service *TourService) Create(tour model.Tour) error {
	if err := service.TourRepository.Save(tour); err != nil {
		return fmt.Errorf("failed to create tour: %w", err)
	}
	return nil
}

func (service *TourService) Delete(id uint64) error {
	if err := service.TourRepository.Delete(id); err != nil {
		return fmt.Errorf("failed to delete tour with ID %d: %w", id, err)
	}
	return nil
}

func (service *TourService) Update(tour model.Tour) error {
	if err := service.TourRepository.Update(tour); err != nil {
		return fmt.Errorf("failed to update tour with ID %d: %w", tour.ID, err)
	}
	return nil
}

func (service *TourService) GetAll() ([]model.Tour, error) {
	tour, err := service.TourRepository.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all tours: %w", err)
	}
	return tour, nil
}

func (service *TourService) GetByID(id uint64) (*model.Tour, error) {
	tour, err := service.TourRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tour with ID %d: %w", id, err)
	}
	return tour, nil
}
