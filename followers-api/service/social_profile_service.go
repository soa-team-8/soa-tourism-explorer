package service

import (
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

func (service *SocialProfileService) Follow(followerID uint64, followedID uint64) error {

	// TODO: implement follow
	return nil
}

func (service *SocialProfileService) Unfollow(followerID uint64, followedID uint64) error {

	// TODO: implement unfollow
	return nil
}
