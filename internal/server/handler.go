package server

import (
	"context"
	"fmt"

	"github.com/fr33dman/go-template/internal/server/gen"
)

type EntityHandler interface {
	CreateEntity(
		ctx context.Context,
		request gen.CreateEntityRequestObject,
	) (gen.CreateEntityResponseObject, error)
	GetEntity(
		ctx context.Context,
		request gen.GetEntityRequestObject,
	) (gen.GetEntityResponseObject, error)
}

type APIHandler struct {
	entityHandler EntityHandler
}

func NewAPIHandler(
	entityHandler EntityHandler,
) *APIHandler {
	return &APIHandler{
		entityHandler: entityHandler,
	}
}

var _ gen.StrictServerInterface = (*APIHandler)(nil)

func (h *APIHandler) CreateEntity(
	ctx context.Context,
	request gen.CreateEntityRequestObject,
) (gen.CreateEntityResponseObject, error) {
	if h.entityHandler == nil {
		return nil, notImplemented("entity handler not implemented")
	}
	return h.entityHandler.CreateEntity(ctx, request)
}

func (h *APIHandler) GetEntity(
	ctx context.Context,
	request gen.GetEntityRequestObject,
) (gen.GetEntityResponseObject, error) {
	if h.entityHandler == nil {
		return nil, notImplemented("entity handler not implemented")
	}
	return h.entityHandler.GetEntity(ctx, request)
}

func notImplemented(name string) error {
	return fmt.Errorf("%s is not configured", name)
}
