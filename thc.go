package thc

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

var removedID = uuid.NewString() // so little truly matters

type dataMap map[string]struct {
	value        any
	timeModified time.Time
}

type container struct {
	identity string
	data     dataMap
	mut      sync.RWMutex // goroutine safety compliance
}

type key[T any] struct {
	identity string
	mapKey   string
}

func (c *container) String() string {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return "Length: " + strconv.Itoa(len(c.data))
}

func (c *container) Len() int {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return len(c.data)
}

func NewTHC() container {
	return container{
		identity: uuid.NewString(),
		data: make(map[string]struct {
			value        any
			timeModified time.Time
		}),
	}
}

func Store[T any](c *container, input T) (key[T], error) {
	switch any(input).(type) {
	case container:
		if any(input).(container).identity == c.identity {
			var zero key[T]
			return zero, fmt.Errorf("container may not store itself")
		}
	}

	newKey := uuid.NewString()

	c.mut.Lock()
	defer c.mut.Unlock()

	c.data[newKey] = struct {
		value        any
		timeModified time.Time
	}{
		value:        input,
		timeModified: time.Now(),
	}

	return key[T]{
		identity: c.identity,
		mapKey:   newKey,
	}, nil
}

func Fetch[T any](c *container, key key[T]) (T, error) {
	var zero T

	if key.identity == removedID {
		return zero, fmt.Errorf("deleted value at key")
	}
	if c.identity != key.identity {
		return zero, fmt.Errorf("container/key identity mismatch")
	}

	c.mut.RLock()
	defer c.mut.RUnlock()

	val, ok := c.data[key.mapKey]
	if !ok {
		return zero, fmt.Errorf("value not found")
	}

	casted, ok := val.value.(T)
	if !ok {
		return zero, fmt.Errorf("type-casting error")
	}
	return casted, nil
}

func Update[T any](c *container, key key[T], input T) error {
	switch any(input).(type) {
	case container:
		if any(input).(container).identity == c.identity {
			return fmt.Errorf("container may not store itself")
		}
	}
	if key.identity == removedID {
		return fmt.Errorf("deleted value at key")
	}
	if c.identity != key.identity {
		return fmt.Errorf("container/key identity mismatch")
	}

	c.mut.Lock()
	defer c.mut.Unlock()

	c.data[key.mapKey] = struct {
		value        any
		timeModified time.Time
	}{
		value:        input,
		timeModified: time.Now(),
	}
	return nil
}

func Remove[T any](c *container, key *key[T]) error {
	if key.identity == removedID {
		return fmt.Errorf("deleted value at key")
	}
	if c.identity != key.identity {
		return fmt.Errorf("container/key identity mismatch")
	}

	c.mut.Lock()
	defer c.mut.Unlock()

	_, ok := c.data[key.mapKey]
	if !ok {
		return fmt.Errorf("no value to remove at key")
	}

	key.identity = removedID
	delete(c.data, key.mapKey)

	return nil
}
