// Package task provides a simple interface for managing tasks or services
// that can be started and stopped. This is useful for long-running processes
// that need to be controlled or managed in a standardized way.
package task

import "context"

// Runner defines an interface for managing the lifecycle of a task or service.
// It provides methods to start the task and to stop it gracefully using a context.
type Runner interface {
	// Start initiates the task or service. It returns an error if the task fails
	// to start or encounters an issue during the startup process.
	Start() error

	// Stop terminates the task or service. The method accepts a context, which
	// can be used to manage timeouts or cancellation signals, ensuring that the
	// task stops gracefully. It returns an error if the task fails to stop cleanly.
	Stop(ctx context.Context) error
}
