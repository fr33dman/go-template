package services

import (
	"context"
	"fmt"

	"github.com/fr33dman/go-template/internal/core"
)

type EntityRepository interface {
	FetchOne(ctx context.Context, id int64) (core.Entity, error)
	Insert(ctx context.Context, entity core.Entity) (core.Entity, error)
}

type Transactor interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type EntityService struct {
	entities EntityRepository
	events   EventPublisher
	tx       Transactor
}

func NewEntityService(
	entities EntityRepository,
	events EventPublisher,
	tx Transactor,
) *EntityService {
	return &EntityService{
		entities: entities,
		events:   events,
		tx:       tx,
	}
}

func (s *EntityService) Get(ctx context.Context, id int64) (core.Entity, error) {
	return s.entities.FetchOne(ctx, id)
}

func (s *EntityService) Create(ctx context.Context, entity core.Entity) (core.Entity, error) {
	var created core.Entity
	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		created, err = s.entities.Insert(ctx, entity)
		if err != nil {
			return err
		}

		if created.Id == nil {
			return fmt.Errorf("entity id is nil after insert")
		}

		return s.events.Publish(ctx, Event{
			AggregateType: core.EntityAggregateType,
			AggregateID:   *created.Id,
			EventType:     core.EntityEventCreated,
			Payload:       created,
		})
	})
	if err != nil {
		return core.Entity{}, err
	}

	return created, nil
}
