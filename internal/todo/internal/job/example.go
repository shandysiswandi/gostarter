package job

import (
	"context"
	"log"
)

type ExampleJob struct{}

func (e *ExampleJob) Start() error {
	log.Println("example job start")

	return nil
}

func (e *ExampleJob) Stop(context.Context) error {
	log.Println("example job stop")

	return nil
}
