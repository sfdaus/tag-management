package http

import (
	"net/http"
	"prakarsa-app/delivery/middleware"
	"prakarsa-app/domain"
	"prakarsa-app/transport/request"
	"prakarsa-app/utils"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

type TagHandler struct {
	TagUC domain.TagUsecase
}

// NewTagHandler will initialize the todo resources endpoint
func NewTagHandler(e *echo.Echo, middleware *middleware.Middleware, tagUC domain.TagUsecase) {
	handler := &TagHandler{
		TagUC: tagUC,
	}

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/tags", handler.Create)
	apiV1.PATCH("/tags/:id", handler.Update)
	apiV1.DELETE("/tags/:id", handler.Delete)
	apiV1.GET("/tags", handler.GetList)
	apiV1.GET("/tags/:id", handler.GetDetail)
}

func (h *TagHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.CreateTagReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if res, err := h.TagUC.Create(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Tag successfully created",
			"data":    res,
		})
	}

}

func (h *TagHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.UpdateTagReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if err := h.TagUC.Update(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Tag successfully updated",
	})
}

func (h *TagHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.DeleteTagReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if rowsAffected, err := h.TagUC.Delete(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Tag successfully deleted",
			"data": map[string]int64{
				"rows_affected": rowsAffected,
			},
		})
	}
}

func (h *TagHandler) GetList(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.GetListTagReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if res, meta, err := h.TagUC.GetList(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Tag successfully retrieved",
			"data": map[string]interface{}{
				"data": res,
				"meta": meta,
			},
		})
	}
}

func (h *TagHandler) GetDetail(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.GetDetailTagReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if res, err := h.TagUC.GetDetail(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Tag successfully retrieved",
			"data": map[string]interface{}{
				"data": res,
			},
		})
	}
}
