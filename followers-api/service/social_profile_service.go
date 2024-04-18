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

func (service *SocialProfileService) Follow(followerID uint64, followedID uint64) error {
	if followerID == followedID {
		return errors.New("you cannot follow yourself")
	}

	// followers social profile
	followerSocialProfile, err := service.SocialProfileRepository.GetByUserID(followerID)
	if err != nil {
		return err
	}
	followerSocialProfile.FollowedIds = append(followerSocialProfile.FollowedIds, &followedID)
	err = service.SocialProfileRepository.Update(followerSocialProfile)
	if err != nil {
		return err
	}

	// followed users social profile
	followedSocialProfile, err := service.SocialProfileRepository.GetByUserID(followedID)
	if err != nil {
		return err
	}
	followedSocialProfile.FollowersIds = append(followedSocialProfile.FollowersIds, &followerID)
	err = service.SocialProfileRepository.Update(followedSocialProfile)
	if err != nil {
		return err
	}

	return nil
}

func (service *SocialProfileService) Unfollow(followerID uint64, followedID uint64) error {
	followedSocialProfile, err := service.SocialProfileRepository.GetByUserID(followedID)
	if err != nil {
		return err
	}
	followedSocialProfile.FollowersIds = removeID(followedSocialProfile.FollowersIds, followerID)
	err = service.SocialProfileRepository.Update(followedSocialProfile)
	if err != nil {
		return err
	}

	followerSocialProfile, err := service.SocialProfileRepository.GetByUserID(followerID)
	if err != nil {
		return err
	}

	followerSocialProfile.FollowedIds = removeID(followerSocialProfile.FollowedIds, followedID)
	err = service.SocialProfileRepository.Update(followerSocialProfile)
	if err != nil {
		return err
	}

	return nil
}

func removeID(ids []*uint64, idToRemove uint64) []*uint64 {
	var updatedIDs []*uint64
	for _, id := range ids {
		if *id != idToRemove {
			updatedIDs = append(updatedIDs, id)
		}
	}
	return updatedIDs
}
