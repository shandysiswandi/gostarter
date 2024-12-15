package usecase

import (
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type Dependency struct {
	UIDNumber uid.NumberID
	Validator validation.Validator
	Telemetry *telemetry.Telemetry
}
