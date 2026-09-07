package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fr33dman/go-template/internal/heartbeat"
	"github.com/fr33dman/go-template/internal/outbox"
	repo "github.com/fr33dman/go-template/internal/repo/entity"
	"github.com/fr33dman/go-template/internal/services"
	"github.com/fr33dman/go-template/pkg/database"
	"github.com/fr33dman/go-template/pkg/probes"
)

type (
	Container struct {
		// DB
		DBPool *pgxpool.Pool

		// Probes
		Probes *probes.Probes

		// Services
		EntityService *services.EntityService
	}
)

var _ services.Transactor = (*database.Transactor)(nil)

func NewContainer(ctx context.Context, cfg Config) (Container, error) {
	// Setup Database Connection
	pool, err := database.GetConnectionPool(ctx, cfg.Database.DSN())
	if err != nil {
		return Container{}, err
	}

	// Setup Probes
	appProbes := heartbeat.NewProbes(pool)

	// Setup Database Adapters
	entityRepository := repo.NewEntityRepository(pool)
	eventPublisher := outbox.NewPublisher()
	transactor := database.NewTransactor(pool)

	return Container{
		DBPool:        pool,
		Probes:        appProbes,
		EntityService: services.NewEntityService(entityRepository, eventPublisher, transactor),
	}, nil
}

func (c Container) Close() {
	if c.DBPool != nil {
		c.DBPool.Close()
	}
}
