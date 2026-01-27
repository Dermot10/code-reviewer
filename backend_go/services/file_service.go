package services

import (
	"fmt"
	"log/slog"

	"github.com/dermot10/code-reviewer/backend_go/models"
	"gorm.io/gorm"
)

type FileService struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewFileService(db *gorm.DB, logger *slog.Logger) *FileService {
	return &FileService{db: db, logger: logger}
}

func (f *FileService) CreateFile(userID uint, name, content string) (*models.File, error) {
	file := &models.File{
		UserID:  userID,
		Name:    name,
		Content: content,
	}

	if err := f.db.Create(file).Error; err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	return file, nil
}

func (f *FileService) UpdateFile(userID, fileID uint, content string) (*models.File, error) {
	var file models.File

	if err := f.db.Where("id = ? AND user_id  = ?", fileID, userID).First(&file).Error; err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	file.Content = content
	if err := f.db.Save(&file).Error; err != nil {
		return nil, fmt.Errorf("faield to update file; %w", err)
	}
	return &file, nil
}

func (f *FileService) GetFile(userID, fileID uint) (*models.File, error) {
	var file models.File
	if err := f.db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	return &file, nil
}

func (f *FileService) ListFiles(userID uint) ([]models.File, error) {
	var files []models.File
	if err := f.db.Where("user_id = ?", userID).Find(&files).Error; err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	return files, nil
}

func (f *FileService) DeleteFile(userID, fileID uint) error {
	result := f.db.Where("id = ? AND user_id = ?", fileID, userID).Delete(&models.File{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete file: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("file not found")
	}
	return nil
}
