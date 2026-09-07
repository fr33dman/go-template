package services_test

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/fr33dman/go-template/internal/core"
	"github.com/fr33dman/go-template/internal/services"
	"github.com/fr33dman/go-template/internal/services/mocks"
)

func TestEntityServiceCreatePublishesCreatedEventInTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	entities := mocks.NewMockEntityRepository(ctrl)
	events := mocks.NewMockEventPublisher(ctrl)
	tx := mocks.NewMockTransactor(ctrl)
	service := services.NewEntityService(entities, events, tx)

	input := core.Entity{Field1: "value", Field2: 42}
	id := int64(10)
	created := core.Entity{Id: &id, Field1: input.Field1, Field2: input.Field2}

	tx.EXPECT().
		WithinTx(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	entities.EXPECT().
		Insert(ctx, input).
		Return(created, nil)
	events.EXPECT().
		Publish(ctx, services.Event{
			AggregateType: core.EntityAggregateType,
			AggregateID:   id,
			EventType:     core.EntityEventCreated,
			Payload:       created,
		}).
		Return(nil)

	got, err := service.Create(ctx, input)
	if err != nil {
		t.Fatalf("create entity: %v", err)
	}

	if got != created {
		t.Fatalf("expected created entity %+v, got %+v", created, got)
	}
}

func TestEntityServiceCreateReturnsInsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	entities := mocks.NewMockEntityRepository(ctrl)
	events := mocks.NewMockEventPublisher(ctrl)
	tx := mocks.NewMockTransactor(ctrl)
	service := services.NewEntityService(entities, events, tx)

	input := core.Entity{Field1: "value", Field2: 42}
	insertErr := errors.New("insert failed")

	tx.EXPECT().
		WithinTx(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	entities.EXPECT().
		Insert(ctx, input).
		Return(core.Entity{}, insertErr)

	_, err := service.Create(ctx, input)
	if !errors.Is(err, insertErr) {
		t.Fatalf("expected insert error, got %v", err)
	}
}

func TestEntityServiceCreateRequiresCreatedID(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	entities := mocks.NewMockEntityRepository(ctrl)
	events := mocks.NewMockEventPublisher(ctrl)
	tx := mocks.NewMockTransactor(ctrl)
	service := services.NewEntityService(entities, events, tx)

	input := core.Entity{Field1: "value", Field2: 42}
	created := core.Entity{Field1: input.Field1, Field2: input.Field2}

	tx.EXPECT().
		WithinTx(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	entities.EXPECT().
		Insert(ctx, input).
		Return(created, nil)

	_, err := service.Create(ctx, input)
	if err == nil {
		t.Fatal("expected error")
	}
}
