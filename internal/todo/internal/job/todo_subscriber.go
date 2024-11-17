package job

import (
	"context"
	"log"

	"github.com/shandysiswandi/gostarter/pkg/messaging"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

type TodoSubscriber struct {
	MsgClient messaging.Client
	Tel       *telemetry.Telemetry
}

func (e *TodoSubscriber) Start() error {
	ctx := context.Background()
	e.Tel.Logger().Debug(ctx, "starting subscription to todo topic")

	fn := func(ctx context.Context, data *messaging.Data) error {
		log.Println("ctx", ctx)
		log.Println("data", data)

		return nil
	}

	return e.MsgClient.Subscribe(ctx, "topic", "subID", fn)
}

func (e *TodoSubscriber) Stop(context.Context) error {
	log.Println("example job stop")

	return nil
}
