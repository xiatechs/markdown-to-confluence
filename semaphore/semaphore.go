// Package semaphore is to control flow of concurrent operations
package semaphore

import (
	"sync"
)

// Semaphore struct is for storing a sync.WaitGroup and a channel
type Semaphore struct {
	c  chan struct{}
	wg *sync.WaitGroup
}

// NewSemaphore returns a Semaphore object
func NewSemaphore(maxConcurrentOps int) *Semaphore {
	return &Semaphore{make(chan struct{}, maxConcurrentOps), new(sync.WaitGroup)}
}

// Add method adds a goroutine to waitgroup & adds to semphore
func (s *Semaphore) Add() {
	s.wg.Add(1)
	s.c <- struct{}{}
}

// Done method removes a goroutine from waitgroup & removes from semphore
func (s *Semaphore) Done() {
	<-s.c
	s.wg.Done()
}

// Wait method is to wait for goroutines to finish
func (s *Semaphore) Wait() {
	s.wg.Wait()
}
