package service

import (
	"errors"
	"followers-api/model"
	repository "followers-api/repository"
)

type SocialProfileService struct {
	SocialProfileRepository *repository.SocialProfileRepository
}

func NewSocialProfileService(repo *repository.SocialProfileRepository) *SocialProfileService {
	return &SocialProfileService{
		SocialProfileRepository: repo,
	}
}

func (service *SocialProfileService) CreateUser(user *model.User) error {
	_, err := service.SocialProfileRepository.SaveUser(user)
	return err
}

func (service *SocialProfileService) Follow(followerID uint64, followedID uint64) (*model.SocialProfile, error) {
	if followerID == followedID {
		return nil, errors.New("you cannot follow yourself")
	}

	err := service.SocialProfileRepository.Follow(followerID, followedID)
	if err != nil {
		return nil, err
	}

	followerSocialProfile, err := service.SocialProfileRepository.GetSocialProfile(followerID)
	if err != nil {
		return nil, err
	}
	return followerSocialProfile, nil
}

func (service *SocialProfileService) Unfollow(followerID uint64, followedID uint64) (*model.SocialProfile, error) {
	err := service.SocialProfileRepository.Unfollow(followerID, followedID)
	if err != nil {
		return nil, err
	}
	followerSocialProfile, err := service.SocialProfileRepository.GetSocialProfile(followerID)
	if err != nil {
		return nil, err
	}
	return followerSocialProfile, nil
}

func (service *SocialProfileService) GetProfile(ID uint64) (*model.SocialProfile, error) {
	socialProfile, err := service.SocialProfileRepository.GetSocialProfile(ID)
	if err != nil {
		return nil, err
	}
	return socialProfile, nil
}

func (service *SocialProfileService) GetRecommendations(ID uint64) ([]*model.User, error) {
	socialProfile, err := service.SocialProfileRepository.GetSocialProfile(ID)
	if err != nil {
		return nil, err
	}
	recommendations, err := service.SocialProfileRepository.GetRecommendations(ID)
	if err != nil {
		return nil, err
	}
	filteredRecommendations := make([]*model.User, 0)
	followedMap := make(map[uint64]bool)
	for _, followed := range socialProfile.Followed {
		followedMap[followed.ID] = true
	}
	for _, recommendation := range recommendations {
		if !followedMap[recommendation.ID] {
			filteredRecommendations = append(filteredRecommendations, recommendation)
		}
	}

	for i, rec := range filteredRecommendations {
		if rec.ID == ID {
			filteredRecommendations = append(filteredRecommendations[:i], filteredRecommendations[i+1:]...)
			break
		}
	}
	return filteredRecommendations, nil
}
