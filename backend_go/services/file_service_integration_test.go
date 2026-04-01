package services

import (
	"testing"

	"github.com/dermot10/code-reviewer/backend_go/models"
	"github.com/stretchr/testify/require"
)

func TestFileService_CreateFile(t *testing.T) {
	setup := SetUpTest(t)
	db := setup.DB

	logger := newTestLogger()
	service := NewFileService(db, logger)

	file, err := service.CreateFile(1, "main.go", "package main")

	require.NoError(t, err)
	require.NotNil(t, file)
	require.Equal(t, "main.go", file.Name)

	var dbFile models.File
	err = db.First(&dbFile, file.ID).Error

	require.NoError(t, err)
	require.Equal(t, "package main", dbFile.Content)
}

func TestFileService_GetFile(t *testing.T) {
	setup := SetUpTest(t)
	db := setup.DB

	logger := newTestLogger()
	service := NewFileService(db, logger)

	file := models.File{
		UserID:  1,
		Name:    "test.go",
		Content: "fmt.Println()",
	}

	require.NoError(t, db.Create(&file).Error)

	result, err := service.GetFile(1, file.ID)

	require.NoError(t, err)
	require.Equal(t, "test.go", result.Name)
}

func TestFileService_UpdateFile(t *testing.T) {
	setup := SetUpTest(t)
	db := setup.DB

	logger := newTestLogger()
	service := NewFileService(db, logger)

	file := models.File{
		UserID:  1,
		Name:    "app.go",
		Content: "old content",
	}

	require.NoError(t, db.Create(&file).Error)

	updated, err := service.UpdateFile(1, file.ID, "new content")

	require.NoError(t, err)
	require.Equal(t, "new content", updated.Content)
}

func TestFileService_ListFiles(t *testing.T) {
	setup := SetUpTest(t)
	db := setup.DB

	logger := newTestLogger()
	service := NewFileService(db, logger)

	files := []models.File{
		{UserID: 1, Name: "a.go", Content: "a"},
		{UserID: 1, Name: "b.go", Content: "b"},
	}

	for _, f := range files {
		require.NoError(t, db.Create(&f).Error)
	}

	result, err := service.ListFiles(1)

	require.NoError(t, err)
	require.Len(t, result, 2)
}

func TestFileService_DeleteFile(t *testing.T) {
	setup := SetUpTest(t)
	db := setup.DB

	logger := newTestLogger()
	service := NewFileService(db, logger)

	file := models.File{
		UserID:  1,
		Name:    "delete.go",
		Content: "test",
	}

	require.NoError(t, db.Create(&file).Error)

	err := service.DeleteFile(1, file.ID)
	require.NoError(t, err)

	var count int64
	db.Model(&models.File{}).Where("id = ?", file.ID).Count(&count)

	require.Equal(t, int64(0), count)
}
