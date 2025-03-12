package job

import (
	"github.com/shandysiswandi/goreng/codec"
	"github.com/shandysiswandi/goreng/config"
	"github.com/shandysiswandi/goreng/messaging"
	"github.com/shandysiswandi/goreng/task"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
)

type Dependency struct {
	Messaging    messaging.Client
	Config       config.Config
	CodecJSON    codec.Codec
	Telemetry    *telemetry.Telemetry
	DomainCreate domain.Create
}

func New(dep Dependency) []task.Runner {
	if !dep.Config.GetBool("feature.flag.todo.job") {
		return nil
	}

	todoSub := &todoSubscriber{
		cjson:        dep.CodecJSON,
		mc:           dep.Messaging,
		tel:          dep.Telemetry,
		createUC:     dep.DomainCreate,
		topic:        "todo.creator.topic",
		subscription: "gostarter.todo.creator.subscription",
	}

	todoPub := &todoPublisher{
		cjson: dep.CodecJSON,
		mc:    dep.Messaging,
		tel:   dep.Telemetry,
		topic: "todo.creator.topic",
	}

	return []task.Runner{
		todoPub,
		todoSub,
	}
}
