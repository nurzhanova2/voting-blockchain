package dto

// CreateElectionRequest — DTO для создания голосования
type CreateElectionRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	IsActive    bool     `json:"is_active"`
	Choices     []string `json:"choices"` 
}

// UpdateElectionRequest — DTO для обновления голосования
type UpdateElectionRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	IsActive    bool     `json:"is_active"`
}
