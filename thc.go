package thc

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

const removedIdentity = "REMOVED"

type thc_container struct {
	identity string
	data     map[string]struct {
		value        any
		timeModified time.Time // i'm thinking so one could sort by time modified
	}
	mut sync.RWMutex // goroutine safety compliance
}

type thc_key[T any] struct {
	identity string
	key      string
}

func (c *thc_container) String() string {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return "Length: " + strconv.Itoa(len(c.data))
}

func (c *thc_container) Len() int {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return len(c.data)
}

func NewTHC() thc_container {
	return thc_container{
		identity: uuid.NewString(),
		data: make(map[string]struct {
			value        any
			timeModified time.Time
		}),
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

	container.mut.Lock()
	defer container.mut.Unlock()

	container.data[key] = struct {
		value        any
		timeModified time.Time
	}{
		value:        input,
		timeModified: time.Now(),
	}

	return thc_key[T]{
		identity: container.identity,
		key:      key,
	}, nil
}

func Fetch[T any](container *thc_container, key thc_key[T]) (T, error) {
	var zero T

	if key.identity == removedIdentity {
		return zero, fmt.Errorf("deleted value at key")
	}
	if container.identity != key.identity {
		return zero, fmt.Errorf("container/key identity mismatch")
	}

	container.mut.RLock()
	defer container.mut.RUnlock()

	val, ok := container.data[key.key]
	if !ok {
		return zero, fmt.Errorf("value not found")
	}

	casted, ok := val.value.(T)
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
	if key.identity == removedIdentity {
		return fmt.Errorf("deleted value at key")
	}
	if container.identity != key.identity {
		return fmt.Errorf("container/key identity mismatch")
	}

	container.mut.Lock()
	defer container.mut.Unlock()

	container.data[key.key] = struct {
		value        any
		timeModified time.Time
	}{
		value:        input,
		timeModified: time.Now(),
	}
	return nil
}

func Remove[T any](container *thc_container, key *thc_key[T]) error {
	if key.identity == removedIdentity {
		return fmt.Errorf("deleted value at key")
	}
	if container.identity != key.identity {
		return fmt.Errorf("container/key identity mismatch")
	}

	container.mut.Lock()
	defer container.mut.Unlock()

	_, ok := container.data[key.key]
	if !ok {
		return fmt.Errorf("no value to remove at key")
	}

	key.identity = removedIdentity
	delete(container.data, key.key)

	return nil
}
