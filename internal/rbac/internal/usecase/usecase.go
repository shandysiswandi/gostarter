package usecase

import (
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/uid"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/pkg/sqlkit"
)

type Dependency struct {
	Telemetry   *telemetry.Telemetry
	Validator   validation.Validator
	UIDNumber   uid.NumberID
	Transaction dbops.Tx
}
