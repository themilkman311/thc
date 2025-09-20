package thc

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

type thc_container struct {
	identity string
	data     map[string]any
}

type thc_key[T any] struct {
	identity string
	key      string
}

func (c thc_container) String() string {
	return "Length: " + strconv.Itoa(len(c.data))
}

func (c thc_container) Len() int {
	return len(c.data)
}

func NewTHC() thc_container {
	return thc_container{
		identity: uuid.NewString(),
		data:     make(map[string]any),
	}
}

func Store[T any](container *thc_container, input T) (thc_key[T], error) {
	switch any(input).(type) {
	case thc_container:
		if any(input).(thc_container).identity == container.identity {
			var zero thc_key[T]
			return zero, fmt.Errorf("container may not store itself")
		}
	}

	key := uuid.NewString()
	container.data[key] = input
	return thc_key[T]{
		identity: container.identity,
		key:      key,
	}, nil
}

func Fetch[T any](container *thc_container, key thc_key[T]) (T, error) {
	var zero T

	if key.identity == "DELETED" {
		return zero, fmt.Errorf("deleted value at key")
	}
	if container.identity != key.identity {
		return zero, fmt.Errorf("container/key identity mismatch")
	}

	val, ok := container.data[key.key]
	if !ok {
		return zero, fmt.Errorf("value not found")
	}

	casted, ok := val.(T)
	if !ok {
		return zero, fmt.Errorf("type-casting error")
	}
	return casted, nil
}

func Update[T any](container *thc_container, key thc_key[T], input T) error {
	switch any(input).(type) {
	case thc_container:
		if any(input).(thc_container).identity == container.identity {
			return fmt.Errorf("container may not store itself")
		}
	}
	if key.identity == "DELETED" {
		return fmt.Errorf("deleted value at key")
	}
	if container.identity != key.identity {
		return fmt.Errorf("container/key identity mismatch")
	}

	container.data[key.key] = input
	return nil
}

func Remove[T any](container *thc_container, key *thc_key[T]) error {
	if key.identity == "DELETED" {
		return fmt.Errorf("deleted value at key")
	}
	if container.identity != key.identity {
		return fmt.Errorf("container/key identity mismatch")
	}
	_, ok := container.data[key.key]
	if !ok {
		return fmt.Errorf("no value to remove at key")
	}

	key.identity = "DELETED"
	delete(container.data, key.key)

	return nil
}
