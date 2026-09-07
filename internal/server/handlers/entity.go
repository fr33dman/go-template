package handlers

import (
	"context"
	"fmt"

	"github.com/fr33dman/go-template/internal/core"
	apperrors "github.com/fr33dman/go-template/internal/errors"
	"github.com/fr33dman/go-template/internal/server/gen"
)

type EntityService interface {
	Get(ctx context.Context, id int64) (core.Entity, error)
	Create(ctx context.Context, entity core.Entity) (core.Entity, error)
}

type EntityHandler struct {
	service EntityService
}

func NewEntityHandler(service EntityService) *EntityHandler {
	return &EntityHandler{service: service}
}

func (h *EntityHandler) CreateEntity(
	ctx context.Context,
	request gen.CreateEntityRequestObject,
) (gen.CreateEntityResponseObject, error) {
	if request.Body == nil {
		return nil, fmt.Errorf("%w: request body is required", apperrors.ErrInvalidArgument)
	}

	newEntity := core.Entity{
		Id:     nil,
		Field1: request.Body.Field1,
		Field2: int(request.Body.Field2),
	}
	entity, err := h.service.Create(ctx, newEntity)
	if err != nil {
		return nil, err
	}
	if entity.Id == nil {
		return nil, fmt.Errorf("entity id is nil")
	}
	return gen.CreateEntity201JSONResponse(mapEntity(entity)), nil
}

func (h *EntityHandler) GetEntity(
	ctx context.Context,
	request gen.GetEntityRequestObject,
) (gen.GetEntityResponseObject, error) {
	entity, err := h.service.Get(ctx, request.EntityId)
	if err != nil {
		return nil, err
	}
	if entity.Id == nil {
		return nil, fmt.Errorf("entity id is nil")
	}
	return gen.GetEntity200JSONResponse(mapEntity(entity)), nil
}
