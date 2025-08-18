package response

type CreateTagRes struct {
	ID string `json:"id"`
}

// Get List Response
type GetListTagRes struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// Get Detail Response
type GetDetailTagRes struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   int64  `json:"created_at"`
	CreatedBy   string `json:"created_by"`
	UpdatedAt   int64  `json:"updated_at"`
	UpdatedBy   string `json:"updated_by"`
}
