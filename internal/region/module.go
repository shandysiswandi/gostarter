package region

import (
	"database/sql"

	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
	inboundhttp "github.com/shandysiswandi/gostarter/internal/region/internal/inbound/http"
	"github.com/shandysiswandi/gostarter/internal/region/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/region/internal/service"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/logger"
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
	Logger    logger.Logger
}

func New(dep Dependency) (*Expose, error) {
	mysqlAdministrative := outbound.NewMysqlRegion(dep.Database)

	searchUC := service.NewSearch(dep.Validator, mysqlAdministrative)

	inboundhttp.RegisterRESTEndpoint(dep.Router, &inboundhttp.Endpoint{
		SearchUC: searchUC,
	})

	return &Expose{}, nil
}
