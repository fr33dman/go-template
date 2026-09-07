package core

const (
	EntityAggregateType = "entity"

	EntityEventCreated = "created"
)

type Entity struct {
	Id     *int64
	Field1 string
	Field2 int
}
