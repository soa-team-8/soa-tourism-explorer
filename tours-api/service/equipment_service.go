package service

import (
	"fmt"
	"tours/model"
	"tours/repository"
)

type EquipmentService struct {
	EquipmentRepository *repository.EquipmentRepository
}

func (service *EquipmentService) Create(equipment model.Equipment) error {
	if err := service.EquipmentRepository.Save(equipment); err != nil {
		return fmt.Errorf("failed to create equipment: %w", err)
	}
	return nil
}

func (service *EquipmentService) Delete(id uint64) error {
	if err := service.EquipmentRepository.Delete(id); err != nil {
		return fmt.Errorf("failed to delete equipment with ID %d: %w", id, err)
	}
	return nil
}

func (service *EquipmentService) Update(equipment model.Equipment) error {
	if err := service.EquipmentRepository.Update(equipment); err != nil {
		return fmt.Errorf("failed to update equipment with ID %d: %w", equipment.ID, err)
	}
	return nil
}

func (service *EquipmentService) GetAll() ([]model.Equipment, error) {
	equipment, err := service.EquipmentRepository.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all equipment: %w", err)
	}
	return equipment, nil
}

func (service *EquipmentService) GetByID(id uint64) (*model.Equipment, error) {
	equipment, err := service.EquipmentRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get equipment with ID %d: %w", id, err)
	}
	return equipment, nil
}
