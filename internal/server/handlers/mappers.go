package handlers

import (
	"github.com/fr33dman/go-template/internal/core"
	"github.com/fr33dman/go-template/internal/server/gen"
)

func mapEntity(source core.Entity) gen.Entity {
	return gen.Entity{
		Id:     *source.Id,
		Field1: source.Field1,
		Field2: int32(source.Field2),
	}
}
