package usecase

import (
	"context"
	"fmt"
	"prakarsa-app/transport/response"
	"time"

	"prakarsa-app/domain"
	"prakarsa-app/repository/redis"
	"prakarsa-app/transport/request"

	"github.com/google/uuid"
)

type TagUsecase struct {
	tagRepo    domain.TagRepository
	redisRepo  redis.RedisRepository
	ctxTimeout time.Duration
}

// NewTagUsecase will create new an tagUsecase object representation of ThreadUsecase interface
func NewTagUsecase(tagRepo domain.TagRepository, redisRepo redis.RedisRepository, ctxTimeout time.Duration) *TagUsecase {
	return &TagUsecase{
		tagRepo:    tagRepo,
		redisRepo:  redisRepo,
		ctxTimeout: ctxTimeout,
	}
}

func (u *TagUsecase) Create(c context.Context, request *request.CreateTagReq) (res response.CreateTagRes, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	// Create Payload
	tagID := uuid.NewString()
	t := true
	tagPayload := &domain.Tag{
		ID:          tagID,
		Name:        request.Name,
		Description: request.Description,
		IsActive:    &t,
		CreatedBy:   "TODO_created_by",
		CreatedAt:   time.Now().Unix(),
	}

	// Response Payload
	res.ID = tagID

	err = u.tagRepo.Create(ctx, tagPayload)
	return
}

func (u *TagUsecase) Update(c context.Context, request *request.UpdateTagReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	// Update Payload
	tagPayload := &domain.Tag{
		ID:        request.ID,
		UpdatedBy: "TODO_updated_by",
		UpdatedAt: time.Now().Unix(),
	}

	if request.Name != "" {
		tagPayload.Name = request.Name
	}
	fmt.Println(request)
	if request.Description != "" {
		tagPayload.Description = request.Description
	}

	if request.IsActive != nil {
		tagPayload.IsActive = request.IsActive
	}

	err = u.tagRepo.Update(ctx, tagPayload)
	return
}
func (u *TagUsecase) Delete(c context.Context, request *request.DeleteTagReq) (rowsAffected int64, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	threadPayload := &domain.Tag{
		ID: request.ID,
	}

	rowsAffected, err = u.tagRepo.Delete(ctx, threadPayload)
	return
}
