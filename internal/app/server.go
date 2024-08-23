// Package app provides a structured framework for initializing, configuring,
// and managing the core components of the application. It handles the setup
// of essential services such as logging, configuration management, database
// connections, Redis client, HTTP router and server, utility libraries, and
// background tasks. The package ensures that all components are properly
// initialized at startup and gracefully shut down when the application stops,
// offering a robust foundation for building and maintaining the application.
package app

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/shandysiswandi/gostarter/pkg/logger"
	"google.golang.org/grpc"
)

// Start initializes and starts the server and listens for termination signals.
// It returns a channel that signals when the application should be terminated.
//
// The function spawns two goroutines:
// 1. One for starting the HTTP server and logging any errors that occur during its execution.
// 2. Another for listening to OS signals (e.g., SIGINT, SIGTERM) and triggering a graceful shutdown.
//
// The returned channel is closed once a termination signal is received and processed.
func (a *App) Start() <-chan struct{} {
	terminateChan := make(chan struct{})

	go func() {
		a.logger.Info(context.TODO(), "http server listening", logger.String("address", a.httpServer.Addr))
		err := a.httpServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("http server:", err)
		}
	}()

	go func() {
		grpcPort := a.config.GetString("server.address.grpc")
		listener, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalln("open tcp listener:", err)
		}

		a.logger.Info(context.TODO(), "grpc server listening", logger.String("address", grpcPort))
		if err := a.grpcServer.Serve(listener); err != nil {
			if !errors.Is(err, grpc.ErrServerStopped) {
				log.Fatalln("grpc server, err:", err)
			}
		}
	}()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		<-sigint

		terminateChan <- struct{}{}
		close(terminateChan)

		log.Println("application gracefully shutdown")
	}()

	return terminateChan
}

// Stop gracefully stops the application by closing any active tasks or jobs
// and releasing any resources held by the application.
//
// It iterates over the list of runnables and closers registered in the application,
// invoking their Stop or Close methods, respectively. If any errors occur during
// the closing of resources, they are logged but do not prevent further closures.
//
// The ctx parameter is used to control the timeout or cancellation of the stop operations.
func (a *App) Stop(ctx context.Context) {
	// close tasks or jobs
	for _, run := range a.runables {
		if err := run.Stop(ctx); err != nil {
			log.Println("failed to close runner", err)
		}
	}

	// close resources
	for _, closer := range a.closersFn {
		if err := closer(ctx); err != nil {
			log.Println("failed to close", err)
		}
	}
}
