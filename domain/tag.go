package domain

import (
	"context"
	"prakarsa-app/transport/request"
)

type Tag struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedBy   string `json:"updated_by"`
	UpdatedAt   int64  `json:"updated_at"`
	DeletedAt   int64  `json:"deleted_at"`
}

// // TagRepository represent the tag repository contract
type TagRepository interface {
	Create(ctx context.Context, tag *Tag) error
	Update(ctx context.Context, tag *Tag) error
	Delete(ctx context.Context, tag *Tag) (int64, error)
}

// TagUsecase represent the tag usecase contract
type TagUsecase interface {
	Create(ctx context.Context, request *request.CreateTagReq) error
	Update(ctx context.Context, request *request.UpdateTagReq) error
	Delete(ctx context.Context, request *request.DeleteTagReq) (int64, error)
}
