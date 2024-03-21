package service

import (
	"fmt"
	"time"
	"tours/model"
	"tours/repository"
)

type ReportedIssueService struct {
	ReportedIssueRepository *repository.ReportedIssueRepository
	TourRepository          *repository.TourRepository
}

func (service *ReportedIssueService) Create(category string, description string, priority uint64, tourID uint64, userID uint64) (model.ReportedIssue, error) {
	reportedIssue := model.ReportedIssue{
		TouristID:   userID,
		TourID:      tourID,
		Category:    category,
		Description: description,
		Priority:    priority,
		Time:        time.Now(),
		Deadline:    nil,
		Closed:      false,
		Resolved:    false,
	}

	newId, err := service.ReportedIssueRepository.Save(reportedIssue)
	if err != nil {
		return model.ReportedIssue{}, fmt.Errorf("failed to create reportedIssue: %w", err)
	}
	reportedIssue.ID = newId
	return reportedIssue, nil
}

func (service *ReportedIssueService) GetAll() ([]model.ReportedIssue, error) {
	reportedIssues, err := service.ReportedIssueRepository.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get reportedIssues: %w", err)
	}
	return reportedIssues, nil
}

func (service *ReportedIssueService) GetByAuthor(authorID uint64) ([]model.ReportedIssue, error) {
	reportedIssues, err := service.ReportedIssueRepository.FindByAuthorId(authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reportedIssues with author ID %d: %w", authorID, err)
	}
	return reportedIssues, nil
}

func (service *ReportedIssueService) GetByTourist(touristID uint64) ([]model.ReportedIssue, error) {
	reportedIssues, err := service.ReportedIssueRepository.FindByTouristId(touristID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reportedIssues with tourist ID %d: %w", touristID, err)
	}
	return reportedIssues, nil
}

func (service *ReportedIssueService) AddComment(id uint64, reportedIssueComment model.ReportedIssueComment) (*model.ReportedIssue, error) {
	reportedIssue, err := service.ReportedIssueRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reportedIssue with ID %d: %w", id, err)
	}
	newComment := model.ReportedIssueComment{
		CreationTime: time.Now(),
		CreatorId:    reportedIssueComment.CreatorId,
		Text:         reportedIssueComment.Text,
	}

	if reportedIssue.Comments == nil {
		reportedIssue.Comments = make([]model.ReportedIssueComment, 0)
	}
	if newComment.Text == "" || newComment.CreatorId == 0 {
		return reportedIssue, fmt.Errorf("comment not valid")
	}
	reportedIssue.Comments = append(reportedIssue.Comments, newComment)

	err = service.ReportedIssueRepository.Update(*reportedIssue)
	if err != nil {
		return &model.ReportedIssue{}, fmt.Errorf("failed to update reportedIssue: %w", err)
	}
	return reportedIssue, nil
}

func (service *ReportedIssueService) AddDeadline(id uint64, date time.Time) (*model.ReportedIssue, error) {
	reportedIssue, err := service.ReportedIssueRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reportedIssue with ID %d: %w", id, err)
	}
	reportedIssue.Deadline = &date
	err = service.ReportedIssueRepository.Update(*reportedIssue)
	if err != nil {
		return &model.ReportedIssue{}, fmt.Errorf("failed to update reportedIssue: %w", err)
	}
	return reportedIssue, nil
}

func (service *ReportedIssueService) PenalizeAuthor(id uint64) (*model.ReportedIssue, error) {
	reportedIssue, err := service.ReportedIssueRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reportedIssue with ID %d: %w", id, err)
	}
	reportedIssue.Closed = true

	tour, err := service.TourRepository.FindByID(reportedIssue.TourID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tour with ID %d: %w", id, err)
	}
	tour.Closed = true

	err = service.TourRepository.Update(*tour)
	if err != nil {
		return &model.ReportedIssue{}, fmt.Errorf("failed to close the tour: %w", err)
	}
	err = service.ReportedIssueRepository.Update(*reportedIssue)
	if err != nil {
		return &model.ReportedIssue{}, fmt.Errorf("failed to update reportedIssue: %w", err)
	}
	return reportedIssue, nil
}

func (service *ReportedIssueService) Close(id uint64) (*model.ReportedIssue, error) {
	reportedIssue, err := service.ReportedIssueRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reportedIssue with ID %d: %w", id, err)
	}
	reportedIssue.Closed = true
	err = service.ReportedIssueRepository.Update(*reportedIssue)
	if err != nil {
		return &model.ReportedIssue{}, fmt.Errorf("failed to update reportedIssue: %w", err)
	}
	return reportedIssue, nil
}

func (service *ReportedIssueService) Resolve(id uint64) (*model.ReportedIssue, error) {
	reportedIssue, err := service.ReportedIssueRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reportedIssue with ID %d: %w", id, err)
	}
	reportedIssue.Resolved = true
	err = service.ReportedIssueRepository.Update(*reportedIssue)
	if err != nil {
		return &model.ReportedIssue{}, fmt.Errorf("failed to update reportedIssue: %w", err)
	}
	return reportedIssue, nil
}
