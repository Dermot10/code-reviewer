package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dermot10/code-reviewer/backend_go/dto"
	"github.com/dermot10/code-reviewer/backend_go/middleware"
	"github.com/dermot10/code-reviewer/backend_go/services"
	"gorm.io/gorm"
)

type FileHandler struct {
	logger      *slog.Logger
	db          *gorm.DB
	fileService *services.FileService
}

func NewFileHandler(logger *slog.Logger, db *gorm.DB, fileService *services.FileService) *FileHandler {
	return &FileHandler{
		logger:      logger,
		db:          db,
		fileService: fileService,
	}
}

func (h *FileHandler) CreateFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, "unathorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "file name required", http.StatusBadRequest)
		return
	}

	file, err := h.fileService.CreateFile(userID, req.Name, req.Content)
	if err != nil {
		h.logger.Error("failed to create file", "error", err)
		http.Error(w, "failed to created file", http.StatusInternalServerError)
		return
	}

	resp := dto.FileResponse{
		ID:        file.ID,
		Name:      file.Name,
		Content:   file.Content,
		CreatedAt: file.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: file.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)

	files, err := h.fileService.ListFiles(userID)
	if err != nil {
		h.logger.Error("failed to list files", "error", err)
		http.Error(w, "failed to list files", http.StatusInternalServerError)
		return
	}

	resp := []dto.FileResponse{}

	for _, file := range files {
		resp = append(resp, dto.FileResponse{
			ID:        file.ID,
			Name:      file.Name,
			Content:   file.Content,
			CreatedAt: file.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: file.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)
	fileID, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "invalid file ID", http.StatusBadRequest)
		return
	}

	file, err := h.fileService.GetFile(userID, uint(fileID))
	if err != nil {
		h.logger.Error("failed to get file", "error", err)
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	resp := dto.FileResponse{
		ID:        file.ID,
		Name:      file.Name,
		Content:   file.Content,
		CreatedAt: file.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: file.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	json.NewEncoder(w).Encode(resp)

}

func (h *FileHandler) UpdateFile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)
	fileID, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "invalid file ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	file, err := h.fileService.UpdateFile(userID, uint(fileID), req.Content)
	if err != nil {
		h.logger.Error("failed to update file", "error", err)
		http.Error(w, "failed to update file", http.StatusInternalServerError)
		return
	}

	resp := dto.FileResponse{
		ID:        file.ID,
		Name:      file.Name,
		Content:   file.Content,
		UpdatedAt: file.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)
	fileID, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "invalid file ID", http.StatusBadRequest)
		return
	}

	if err := h.fileService.DeleteFile(userID, uint(fileID)); err != nil {
		h.logger.Error("failed to delete file", "error", err)
		http.Error(w, "failed to delete file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
