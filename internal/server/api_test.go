package server

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fr33dman/go-template/internal/core"
	apperrors "github.com/fr33dman/go-template/internal/errors"
	"github.com/fr33dman/go-template/internal/server/handlers"
)

type fakeEntityService struct{}

var _ handlers.EntityService = (*fakeEntityService)(nil)

func (s *fakeEntityService) Get(ctx context.Context, id int64) (core.Entity, error) {
	return core.Entity{}, apperrors.ErrNotFound
}

func (s *fakeEntityService) Create(ctx context.Context, entity core.Entity) (core.Entity, error) {
	id := int64(1)
	entity.Id = &id
	return entity, nil
}

func TestAPIServerValidatesRequests(t *testing.T) {
	apiServer := NewAPIServer("127.0.0.1", 0, testLogger())
	err := apiServer.RegisterHandlers(handlers.NewEntityHandler(&fakeEntityService{}))
	if err != nil {
		t.Fatalf("register handlers: %v", err)
	}

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/api/v1/entities",
		strings.NewReader(`{"field1":"value","field2":42,"extra":"rejected"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	apiServer.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusBadRequest,
			rec.Code,
			rec.Body.String(),
		)
	}
}

func TestAPIServerMapsNotFoundError(t *testing.T) {
	apiServer := NewAPIServer("127.0.0.1", 0, testLogger())
	err := apiServer.RegisterHandlers(handlers.NewEntityHandler(&fakeEntityService{}))
	if err != nil {
		t.Fatalf("register handlers: %v", err)
	}

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/v1/entities/1",
		nil,
	)
	rec := httptest.NewRecorder()

	apiServer.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}
