package dto

type CreateFileRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type UpdateFileRequest struct {
	Content string `json:"content"`
}

type FileResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
