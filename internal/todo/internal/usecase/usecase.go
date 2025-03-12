package usecase

import (
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/uid"
	"github.com/shandysiswandi/goreng/validation"
)

type Dependency struct {
	UIDNumber uid.NumberID
	Validator validation.Validator
	Telemetry *telemetry.Telemetry
}
