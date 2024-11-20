package job

import (
	"context"
	"fmt"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/messaging"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

type todoPublisher struct {
	cjson codec.Codec
	mc    messaging.Client
	tel   *telemetry.Telemetry
	topic string
}

func (e *todoPublisher) Start() error {
	ctx, span := e.tel.Tracer().Start(context.Background(), "job.todoPublisher.Start")
	defer span.End()

	e.tel.Logger().Info(ctx, "todo publisher has started")

	go func() {
		messages := make([]*messaging.Data, 0, 10)
		for i := range 10 {
			index := i + 1
			todo := domain.Todo{
				Title:       fmt.Sprintf("title from publisher %d", index),
				Description: fmt.Sprintf("description from publisher %d, this only for example", index),
			}

			bt, err := e.cjson.Encode(todo)
			if err != nil {
				e.tel.Logger().Error(ctx, "failed encode to json", err)

				return
			}

			messages = append(messages, &messaging.Data{Msg: bt, Attributes: nil})
		}

		if err := e.mc.Publish(ctx, e.topic, messages[0]); err != nil {
			e.tel.Logger().Error(ctx, "failed publish message", err)

			return
		}

		if err := e.mc.BulkPublish(ctx, e.topic, messages[1:]); err != nil {
			e.tel.Logger().Error(ctx, "failed bulk publish message", err)

			return
		}
	}()

	return nil
}

func (e *todoPublisher) Stop(ctx context.Context) error {
	e.tel.Logger().Info(ctx, "todo publisher has stopped")

	return nil
}
