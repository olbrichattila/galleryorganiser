package main

import (
	"sync"
)

type SaveCounter interface {
	Increment()
	Decrement()
	Value() int
}

type counter struct {
	mu sync.Mutex
	v  int
}

func (c *counter) Increment() {
	c.mu.Lock()
	c.v++
	c.mu.Unlock()
}

func (c *counter) Decrement() {
	c.mu.Lock()
	c.v--
	c.mu.Unlock()
}

func (c *counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v
}
