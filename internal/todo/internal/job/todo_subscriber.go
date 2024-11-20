package job

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/messaging"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

type todoSubscriber struct {
	cjson               codec.Codec
	mc                  messaging.Client
	tel                 *telemetry.Telemetry
	createUC            domain.Create
	topic, subscription string
}

func (e *todoSubscriber) Start() error {
	ctx := context.Background()
	e.tel.Logger().Info(ctx, "todo subscriber has started")

	return e.mc.Subscribe(ctx, e.topic, e.subscription, e.do)
}

func (e *todoSubscriber) do(ctx context.Context, data *messaging.Data) error {
	var todo domain.Todo
	if err := e.cjson.Decode(data.Msg, &todo); err != nil {
		return err
	}

	_, err := e.createUC.Call(ctx, domain.CreateInput{
		Title:       todo.Title,
		Description: todo.Description,
	})

	return err
}

func (e *todoSubscriber) Stop(ctx context.Context) error {
	e.tel.Logger().Info(ctx, "todo subscriber has stopped")

	return nil
}
