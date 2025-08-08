package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// CreateTagReq represent create request body
type CreateTagReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (request CreateTagReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
	)
}

// Update request body
type UpdateTagReq struct {
	ID          string `param:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

func (request UpdateTagReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
	)
}

// Delete request body
type DeleteTagReq struct {
	ID string `param:"id"`
}

func (request DeleteTagReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
	)
}

// GetList request body
type GetListTagReq struct {
	Name     string `query:"name"`
	IsActive *bool  `query:"is_active"`
	PerPage  int64  `query:"per_page"`
	Page     int64  `query:"page"`
}

func (request GetListTagReq) Validate() error {
	return validation.ValidateStruct(
		&request,
	)
}

// GetDetail request body
type GetDetailTagReq struct {
	ID string `param:"id"`
}

func (request GetDetailTagReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
	)
}
