package usecase

import (
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type Dependency struct {
	Telemetry   *telemetry.Telemetry
	Validator   validation.Validator
	UIDNumber   uid.NumberID
	Transaction dbops.Tx
}
