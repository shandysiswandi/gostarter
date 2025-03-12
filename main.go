// Package main serves as the entry point for the application.
// It initializes the application, starts it, and ensures a graceful shutdown
// by listening for termination signals.
package main

import (
	"context"
	"time"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/shandysiswandi/gostarter/internal/app"
)

// main is the entry point of the application.
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	application := app.New()    // Initialize the application
	wait := application.Start() // Start the application and wait for the termination signal
	<-wait                      // Wait for the application to receive a termination signal
	application.Stop(ctx)       // Stop the application gracefully
}
