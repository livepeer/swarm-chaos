package model

type EntityType int
type OperationType int
type StatusType int

const (
	EntityTypeContainer EntityType = iota
	EntityTypeProcess
	EntityTypeNetworkLink
	EntityTypeVM

	OperationTypeDestroy OperationType = iota
	OperationTypeStart
	OperationTypeStop
	OperationTypePause
	OperationTypeResume
	OperationTypeSlowdown

	StatusTypeWorking StatusType = iota
	StatusTypeDestroyed
	StatusTypeStopped
	StatusTypePaused
	StatusTypeSlow
)

type (
	// Entity represents singles object manageable by Swarm Chaos
	Entity interface {
		Name() string
		Labels() map[string]string
		Childs() []Entity
		Type() EntityType
		Do(operation OperationType) error
		Status() (StatusType, error)
	}

	// Playground represents all the entities that Swarm Chaos can work with.
	Playground interface {
		Entities() ([]Entity, error)
		// EntitiesByLabel(key, value string) ([]Entity, error)
	}

	// Task represents taks
	// Task interface {
	// 	Start() (chan interface{}, error)
	// }
)

// SwarmChaosVersion version
// content of this constant will be set at build time,
// using -ldflags, combining content of `VERSION` file and
// output of the `git describe` command.
var SwarmChaosVersion = "undefined"
