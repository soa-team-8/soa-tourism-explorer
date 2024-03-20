package service

import (
	"fmt"
	"tours/model"
	"tours/repository"
)

type TourService struct {
	TourRepository       *repository.TourRepository
	EquipmentRepository  *repository.EquipmentRepository
	TourRatingRepository *repository.TourRatingRepository
}

func (service *TourService) Create(tour model.Tour) (uint64, error) {
	id, err := service.TourRepository.Save(tour)
	if err != nil {
		return 0, fmt.Errorf("failed to create tour: %w", err)
	}
	return id, nil
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

func (service *TourService) AddEquipmentToTour(tourID, equipmentID uint64) error {
	err := service.TourRepository.AddEquipmentToTour(tourID, equipmentID)
	if err != nil {
		return err
	}

	return nil
}

func (service *TourService) RemoveEquipmentFromTour(tourID, equipmentID uint64) error {
	err := service.TourRepository.RemoveEquipmentFromTour(tourID, equipmentID)
	if err != nil {
		return err
	}

	return nil
}

func (service *TourService) GetToursByAuthor(authorID uint64) ([]model.Tour, error) {
	tours, err := service.TourRepository.FindByAuthorID(authorID)
	if err != nil {
		return nil, err
	}
	return tours, nil
}

func (service *TourService) GetPublishedTours() ([]model.TourPreview, error) {
	tours, err := service.TourRepository.GetPublishedTours()
	if err != nil {
		return nil, fmt.Errorf("failed to get all tours: %w", err)
	}

	var tourPreviews []model.TourPreview
	for _, tour := range tours {
		tourPreview := model.TourPreview{
			ID:                tour.ID,
			AuthorID:          tour.AuthorID,
			Name:              tour.Name,
			Description:       tour.Description,
			DemandignessLevel: tour.DemandignessLevel,
			Price:             tour.Price,
			Tags:              tour.Tags,
			Equipment:         tour.Equipment,
		}

		if len(tour.Checkpoints) > 0 {
			tourPreview.Checkpoint = tour.Checkpoints[0]
		}

		tourPreview.TourRatings, err = service.TourRatingRepository.FindByTourID(tour.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get all tourRatings: %w", err)
		}
		tourPreviews = append(tourPreviews, tourPreview)
	}

	return tourPreviews, nil
}

func (service *TourService) GetPublishedTour(id uint64) (*model.TourPreview, error) {
	tour, err := service.TourRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tour: %w", err)
	}

	tourPreview := model.TourPreview{
		ID:                tour.ID,
		AuthorID:          tour.AuthorID,
		Name:              tour.Name,
		Description:       tour.Description,
		DemandignessLevel: tour.DemandignessLevel,
		Price:             tour.Price,
		Tags:              tour.Tags,
		Equipment:         tour.Equipment,
	}

	if len(tour.Checkpoints) > 0 {
		tourPreview.Checkpoint = tour.Checkpoints[0]
	}

	tourPreview.TourRatings, err = service.TourRatingRepository.FindByTourID(tour.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tourRatings: %w", err)
	}

	return &tourPreview, nil
}

func (service *TourService) GetAverageRating(id uint64) (float32, error) {
	tour, err := service.GetPublishedTour(id)
	if err != nil {
		return 0, fmt.Errorf("failed to get tour: %w", err)
	}

	if len(tour.TourRatings) == 0 {
		return 0, nil
	}

	var totalRating uint64
	for _, rating := range tour.TourRatings {
		totalRating += rating.Rating
	}

	average := float32(totalRating) / float32(len(tour.TourRatings))
	return average, nil
}
