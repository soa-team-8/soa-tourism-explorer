package service

import (
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type ImageService struct {
	ImageStoragePath string
}

func NewImageService() *ImageService {
	imageService := &ImageService{
		ImageStoragePath: "wwwroot/images",
	}
	err := imageService.initializeStorageDirectory()
	if err != nil {
		return nil
	}
	return imageService
}

func (is *ImageService) UploadImages(images []*multipart.FileHeader) ([]string, error) {
	uploadedImageNames := make([]string, 0)

	for _, image := range images {
		file, err := image.Open()
		if err != nil {
			return nil, err
		}

		imageName, err := is.saveImage(file)
		if err != nil {
			err := file.Close()
			if err != nil {
				return nil, err
			}
			return nil, err
		}

		err = file.Close()
		if err != nil {
			return nil, err
		}

		uploadedImageNames = append(uploadedImageNames, imageName)
	}

	return uploadedImageNames, nil
}

func (is *ImageService) saveImage(image multipart.File) (string, error) {
	imageBytes, err := io.ReadAll(image)
	if err != nil {
		return "", err
	}

	imageName := generateUniqueImageName() + ".jpg"
	imagePath := filepath.Join(is.ImageStoragePath, imageName)

	err = os.WriteFile(imagePath, imageBytes, 0644)
	if err != nil {
		return "", err
	}

	return imageName, nil
}

func (is *ImageService) initializeStorageDirectory() error {
	if _, err := os.Stat(is.ImageStoragePath); os.IsNotExist(err) {
		err := os.MkdirAll(is.ImageStoragePath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateUniqueImageName() string {
	return strings.ReplaceAll(
		uuid.New().String(),
		"-",
		"",
	)
}
