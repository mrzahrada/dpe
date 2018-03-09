package store

import (
	"errors"

	"github.com/google/uuid"
	"github.com/mrzahrada/dpe/pkg/fsm"
)

// InMemoryStore implements in memory storage
type InMemoryStore struct {
	StateMachines map[string]fsm.StateMachine
	Executions    map[string]*fsm.Execution
	Wait          map[uint64]string
}

// AddFsm new statemachine
func (store *InMemoryStore) AddFsm(uuid uuid.UUID, fsm fsm.StateMachine) error {
	// check if exists
	if _, ok := store.StateMachines[uuid.String()]; ok {
		return errors.New("state machine already exists")
	}
	store.StateMachines[uuid.String()] = fsm
	return nil
}

// AddExec ...
func (store *InMemoryStore) AddExec(uuid uuid.UUID, exec *fsm.Execution) error {
	// check if exists
	if _, ok := store.Executions[uuid.String()]; ok {
		return errors.New("execution already exists")
	}
	store.Executions[uuid.String()] = exec
	return nil
}

// DeleteExec ...
func (store *InMemoryStore) DeleteExec(uuid uuid.UUID) error {
	delete(store.Executions, uuid.String())
	return nil
}

// UpdateExec ...
func (store *InMemoryStore) UpdateExec(uuid uuid.UUID, state string) error {
	return errors.New("not implemented")
}
