package service

import (
	"errors"
	"followers/model"
	repository "followers/repository"
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
	followerSocialProfile, err := service.SocialProfileRepository.GetSocialProfile(ID)
	if err != nil {
		return nil, err
	}
	return followerSocialProfile, nil
}
