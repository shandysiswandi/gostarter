package usecase

import (
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/goroutine"
	"github.com/shandysiswandi/gostarter/pkg/messaging"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
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
}
