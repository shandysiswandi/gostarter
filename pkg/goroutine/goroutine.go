/*
Package goroutine provides utilities for managing and controlling concurrent goroutines.
It simplifies the process of executing and coordinating multiple tasks concurrently,
with built-in synchronization, error handling, and limiting the maximum number of active goroutines
at any given time. This helps prevent resource exhaustion and ensures smooth execution of concurrent operations.

Example usage:

	package main

	import (
	    "context"
	    "log"
	    "time"

	    ".../goroutine"
	)

	func main() {
	    g := goroutine.NewManager(3)

	    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	    defer cancel()

	    g.Go(ctx, func(_ context.Context) error {
	        log.Println("Task 1 started")
	        time.Sleep(1 * time.Second) // Simulate work
	        log.Println("Task 1 finished")
	        return nil
	    })
	    g.Go(ctx, func(_ context.Context) error {
	        log.Println("Task 2 started")
	        time.Sleep(7 * time.Second) // Simulate not work, context timeout
	        log.Println("Task 2 finished")
	        return nil
	    })
	    g.Go(ctx, func(_ context.Context) error {
	        log.Println("Task 3 started")
	        time.Sleep(5 * time.Second) // Simulate work
	        log.Println("Task 3 finished")
	        return nil
	    })

	    log.Println("Waiting for all tasks to finish")
	    if err := g.Wait(ctx); err != nil {
	        log.Println("Errors occurred:", err)
	    }
	}
*/
package goroutine

import (
	"context"
	"errors"
	"log"
	"runtime/debug"
	"sync"
)

// Manager is type provides methods to start goroutines and wait for their completion.
type Manager struct {
	mu   sync.Mutex
	errs []error
	wg   *sync.WaitGroup
	sema chan struct{}
}

// NewManager initializes a new Manager with a limit on the number of concurrent goroutines.
func NewManager(maxGoroutine int) *Manager {
	return &Manager{
		wg:   &sync.WaitGroup{},
		sema: make(chan struct{}, maxGoroutine), // Semaphore to limit goroutines
	}
}

// Go starts a new goroutine and limits the number of concurrent goroutines with
// a semaphore. If the context is canceled, the goroutine will stop processing
// and will log the cancellation.
func (g *Manager) Go(ctx context.Context, f func(c context.Context) error) {
	select {
	case g.sema <- struct{}{}: // Acquire a semaphore slot
		g.wg.Add(1)

		go func(gCtx context.Context) {
			defer func() {
				<-g.sema // Release semaphore slot
				g.wg.Done()

				if rvr := recover(); rvr != nil {
					log.Println("panic occurred in goroutine")
					debug.PrintStack()
				}
			}()

			select {
			case <-gCtx.Done():
				log.Println("goroutine canceled due to:", gCtx.Err())
			default:
				if err := f(gCtx); err != nil {
					g.mu.Lock()
					g.errs = append(g.errs, err)
					g.mu.Unlock()
				}
			}
		}(ctx)

	default:
		log.Println("Maximum goroutine limit reached, failed to start new goroutine")
	}
}

// Wait waits for all goroutines to complete or for the provided context to be canceled.
// It returns an aggregated error if any goroutines failed or context is timeout/canceled.
func (g *Manager) Wait(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		defer close(done)
		g.wg.Wait()
	}()

	select {
	case <-done:
		return errors.Join(g.errs...)
	case <-ctx.Done():
		return ctx.Err()
	}
}
