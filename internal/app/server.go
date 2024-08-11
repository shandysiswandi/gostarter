package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func (a *App) Start() <-chan struct{} {
	a.ensureInitialized()
	a.ensureClosed()

	terminateChan := make(chan struct{})

	go func() {
		log.Println("http server listen on", a.httpServer.Addr)
		err := a.httpServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("http server:", err)
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

func (a *App) Stop(ctx context.Context) {
	for _, closer := range a.closersFn {
		if err := closer(ctx); err != nil {
			log.Println("failed to close", err)
		}
	}
}
