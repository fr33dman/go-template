package main

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/fr33dman/gobox"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/fr33dman/go-template/internal/app"
	"github.com/fr33dman/go-template/internal/build"
	"github.com/fr33dman/go-template/internal/logger"
	"github.com/fr33dman/go-template/sql/migrations"
)

func main() {
	ctx := context.Background()
	appLogger, _ := logger.New("", build.Version)
	slog.SetDefault(appLogger)

	cfg, err := app.NewConfig()
	if err != nil {
		logger.Fatal(ctx, appLogger, "load config", err)
	}

	appLogger, err = logger.New(cfg.Logger.Level, build.Version)
	if err != nil {
		logger.Fatal(ctx, appLogger, "create logger", err)
	}
	slog.SetDefault(appLogger)

	// Connect to Database
	db, err := sql.Open("postgres", cfg.Database.DSN())
	if err != nil {
		logger.Fatal(ctx, appLogger, "open database", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		logger.Fatal(ctx, appLogger, "ping database", err)
	}

	// Make application migrations
	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Fatal(ctx, appLogger, "set goose dialect", err)
	}

	if err := goose.Up(db, migrations.Dir); err != nil {
		logger.Fatal(ctx, appLogger, "run migrations", err)
	}

	// Make outbox library migrations
	outboxConfig := gobox.OutboxConfiguration{}
	err = outboxConfig.LoadFromEnv()
	if err != nil {
		logger.Fatal(ctx, appLogger, "load outbox configuration", err)
	}
	err = gobox.MigrateOutboxTable(ctx, db, outboxConfig)
	if err != nil {
		logger.Fatal(ctx, appLogger, "run outbox migrations", err)
	}
}
