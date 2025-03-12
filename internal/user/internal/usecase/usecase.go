package usecase

import (
	"github.com/shandysiswandi/goreng/hash"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/validation"
)

type Dependency struct {
	Telemetry *telemetry.Telemetry
	Validator validation.Validator
	Hash      hash.Hash
}
