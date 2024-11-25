package usecase

import (
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type Dependency struct {
	Telemetry *telemetry.Telemetry
	Validator validation.Validator
}
