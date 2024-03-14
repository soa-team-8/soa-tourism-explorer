package repository

//type CheckpointCompletionRepository struct {
//	DB *gorm.DB
//}
//
//func (repo *TourRatingRepository) Save(tourRating model.TourRating) error {
//	result := repo.DB.Create(&tourRating)
//	if result.Error != nil {
//		return result.Error
//	}
//	return nil
//}
//
//func (repo *TourRatingRepository) Delete(id uint64) error {
//	result := repo.DB.Delete(&model.TourRating{}, id)
//	if result.Error != nil {
//		return result.Error
//	}
//	if result.RowsAffected == 0 {
//		return gorm.ErrRecordNotFound
//	}
//	return nil
//}
//
//func (repo *TourRatingRepository) Update(tourRating model.TourRating) error {
//	existingTourRating, err := repo.FindByID(tourRating.ID)
//	if err != nil {
//		return err
//	}
//	if existingTourRating == nil {
//		return errors.New("tourRating not found")
//	}
//
//	result := repo.DB.Save(&tourRating)
//	if result.Error != nil {
//		return result.Error
//	}
//	return nil
//}
//
//func (repo *TourRatingRepository) FindAll() ([]model.TourRating, error) {
//	var tourRating []model.TourRating
//	if err := repo.DB.Find(&tourRating).Error; err != nil {
//		return nil, err
//	}
//	return tourRating, nil
//}
//
//func (repo *TourRatingRepository) FindAllPaged(page, pageSize int) ([]model.TourRating, error) {
//	var tourRating []model.TourRating
//	offset := (page - 1) * pageSize
//	result := repo.DB.Offset(offset).Limit(pageSize).Find(&tourRating)
//	if result.Error != nil {
//		return nil, result.Error
//	}
//	return tourRating, nil
//}
//
//func (repo *TourRatingRepository) FindByID(id uint64) (*model.TourRating, error) {
//	var tourRating model.TourRating
//	if err := repo.DB.First(&tourRating, id).Error; err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return nil, nil
//		}
//		return nil, err
//	}
//	return &tourRating, nil
//}
