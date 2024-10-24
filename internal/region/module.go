package region

import (
	"database/sql"

	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
	"github.com/shandysiswandi/gostarter/internal/region/internal/inbound/http"
	"github.com/shandysiswandi/gostarter/internal/region/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/region/internal/service"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type Expose struct{}

type Dependency struct {
	Database  *sql.DB
	RedisDB   *redis.Client
	Config    config.Config
	CodecJSON codec.Codec
	Validator validation.Validator
	Router    *httprouter.Router
	Telemetry *telemetry.Telemetry
}

func New(dep Dependency) (*Expose, error) {
	sqlRegion := outbound.NewSQLRegion(dep.Database, dep.Config)

	searchService := service.NewSearch(dep.Validator, sqlRegion)

	http.RegisterRESTEndpoint(dep.Router, http.NewEndpoint(searchService))

	return &Expose{}, nil
}
