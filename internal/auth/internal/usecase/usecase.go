package usecase

import (
	"github.com/shandysiswandi/goreng/clock"
	"github.com/shandysiswandi/goreng/hash"
	"github.com/shandysiswandi/goreng/jwt"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/uid"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/pkg/sqlkit"
)

type Dependency struct {
	Telemetry   *telemetry.Telemetry
	Validator   validation.Validator
	UIDNumber   uid.NumberID
	Hash        hash.Hash
	SecHash     hash.Hash
	JWT         jwt.JWT
	Clock       clock.Clocker
	Transaction sqlkit.Tx
}
