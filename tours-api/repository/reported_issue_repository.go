package repository

import (
	"errors"
	"gorm.io/gorm"
	"tours/model"
)

type ReportedIssueRepository struct {
	DB *gorm.DB
}

func (repo *ReportedIssueRepository) Save(reportedIssue model.ReportedIssue) (uint64, error) {
	result := repo.DB.Create(&reportedIssue)
	if result.Error != nil {
		return 0, result.Error
	}
	return reportedIssue.ID, nil
}

func (repo *ReportedIssueRepository) Update(reportedIssue model.ReportedIssue) error {
	existingReportedIssue, err := repo.FindByID(reportedIssue.ID)
	if err != nil {
		return err
	}
	if existingReportedIssue == nil {
		return errors.New("reportedIssue not found")
	}

	result := repo.DB.Save(&reportedIssue)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ReportedIssueRepository) FindByID(id uint64) (*model.ReportedIssue, error) {
	var reportedIssue model.ReportedIssue
	if err := repo.DB.Preload("Comments").Preload("Tour").First(&reportedIssue, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reportedIssue, nil
}

func (repo *ReportedIssueRepository) FindAll() ([]model.ReportedIssue, error) {
	var reportedIssues []model.ReportedIssue
	if err := repo.DB.Preload("Comments").Preload("Tour").Find(&reportedIssues).Error; err != nil {
		return nil, err
	}
	return reportedIssues, nil
}

func (repo *ReportedIssueRepository) FindByAuthorId(authorId uint64) ([]model.ReportedIssue, error) {
	var reportedIssues []model.ReportedIssue
	if err := repo.DB.
		Preload("Comments").Preload("Tour").
		Joins("JOIN tours ON tours.id = reported_issues.tour_id").
		Where("tours.author_id = ?", authorId).
		Find(&reportedIssues).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return reportedIssues, nil
}

func (repo *ReportedIssueRepository) FindByTouristId(touristId uint64) ([]model.ReportedIssue, error) {
	var reportedIssues []model.ReportedIssue
	if err := repo.DB.Preload("Comments").Preload("Tour").Where("tourist_id = ?", touristId).Find(&reportedIssues).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return reportedIssues, nil
}
