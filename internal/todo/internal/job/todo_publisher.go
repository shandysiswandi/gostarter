package job

import (
	"context"
	"log"
	"time"

	"github.com/shandysiswandi/gostarter/pkg/messaging"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

type TodoPublisher struct {
	MsgClient messaging.Client
	Tel       *telemetry.Telemetry
}

func (e *TodoPublisher) Start() error {
	ctx := context.Background()
	e.Tel.Logger().Debug(ctx, "starting publish to todo topic")

	go func() {
		for i := 0; i < 10; i++ {
			err := e.MsgClient.Publish(ctx, "topic", &messaging.Data{Msg: []byte(`{"msg":"hello world"}`)})
			if err != nil {
				return
			}

			time.Sleep(time.Second * 5)
		}

		err := e.MsgClient.BulkPublish(ctx, "topic", []*messaging.Data{
			{
				Msg:        []byte(`{"msg":"hello world 1"}`),
				Attributes: nil,
			},
			{
				Msg:        []byte(`{"msg":"hello world 2"}`),
				Attributes: nil,
			},
			{
				Msg:        []byte(`{"msg":"hello world  3"}`),
				Attributes: nil,
			},
			{
				Msg:        []byte(`{"msg":"hello world 4"}`),
				Attributes: nil,
			},
			{
				Msg:        []byte(`{"msg":"hello world 5"}`),
				Attributes: nil,
			},
		})
		if err != nil {
			return
		}

	}()

	return nil
}

func (e *TodoPublisher) Stop(context.Context) error {
	log.Println("example job stop")

	return nil
}
