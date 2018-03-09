package fsm

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

// Execution defines current state of fsm
type Execution struct {
	uuid uuid.UUID
	data []byte
}

// ApplyPath applies input/result path for Execution
func (state *Execution) ApplyPath(path string) error {
	return errors.New("not implemented")
}

// NewState initialize state
func NewState(input []byte) *Execution {
	return &Execution{
		uuid: uuid.New(),
		data: input,
	}
}

// State abstracts any state, which can be executed
type State interface {
	Execute(*Execution) (string, error)
	// gob marshaling
}

// PassState ...
type PassState struct {
	Next string
	// Result     interface{}
	// ResultPath string
}

// Execute implementes Executor interface
func (p PassState) Execute(_ *Execution) (string, error) {
	log.Println("[DEBUG] Executing pass state")
	return p.Next, nil
}

// WaitDelayState ...
type WaitDelayState struct {
	Delay time.Duration
	Next  string
}

// Execute implementes Executor interface
func (wd WaitDelayState) Execute(_ *Execution) (string, error) {
	log.Println("[DEBUG] Executing wait state with delay", wd.Delay, "seconds")
	time.Sleep(wd.Delay)
	return wd.Next, nil
}

// ChoiceState ...
type ChoiceState struct {
	Choices []Choice
	Next    []string
	Default string
}

// Execute implementes Executor interface
func (ch ChoiceState) Execute(exec *Execution) (string, error) {
	for i, choice := range ch.Choices {
		result, err := choice.Evaluate(exec)
		if err != nil {
			return "", err
		}
		if result {
			return ch.Next[i], nil
		}
	}
	return ch.Default, nil
}

// StateMachine ...
type StateMachine struct {
	States       map[string]State
	StartAt      string
	MaxExecution int
}

// Execute implementes Executor interface
func (sm StateMachine) Execute(exec *Execution) (string, error) {
	state := sm.StartAt
	for i := 0; i < sm.MaxExecution; i++ {
		state, err := sm.States[state].Execute(exec)
		if err != nil {
			return "", err
		}
		if state == "EndState" {
			break
		}
	}
	return "EndState", nil
}
