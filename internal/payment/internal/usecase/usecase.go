package usecase

import (
	"github.com/shandysiswandi/goreng/clock"
	"github.com/shandysiswandi/goreng/codec"
	"github.com/shandysiswandi/goreng/config"
	"github.com/shandysiswandi/goreng/goroutine"
	"github.com/shandysiswandi/goreng/messaging"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/uid"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/pkg/sqlkit"
)

type Dependency struct {
	Messaging   messaging.Client
	Config      config.Config
	UIDNumber   uid.NumberID
	CodecJSON   codec.Codec
	Validator   validation.Validator
	Transaction dbops.Tx
	Telemetry   *telemetry.Telemetry
	Goroutine   *goroutine.Manager
	Clock       clock.Clocker
}
