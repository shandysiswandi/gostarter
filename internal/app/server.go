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
		log.Println("http server listening", "address", a.httpServer.Addr)
		err := a.httpServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("http server:", err)
		}
	}()

	go func() {
		log.Println("gql server listening", "address", a.gqlServer.Addr)
		err := a.gqlServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("gql server:", err)
		}
	}()

	go func() {
		grpcPort := a.config.GetString("server.address.grpc")
		listener, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalln("open tcp listener:", err)
		}

		log.Println("grpc server listening", "address", grpcPort)
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
	for _, run := range a.runnables {
		if err := run.Stop(ctx); err != nil {
			log.Println("failed to close runner", err)
		}
	}

	// close resources
	for name, closer := range a.closerFn {
		if err := closer(ctx); err != nil {
			log.Printf("failed to close %s because: %v", name, err)
		}
	}

	log.Println("waiting for all goroutine to finish")
	if err := a.goroutine.Wait(); err != nil {
		log.Println("error from goroutines executions:", err)
	}
	log.Println("all goroutines have finished successfully")
}
