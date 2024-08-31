package shortly

import (
	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
	"github.com/shandysiswandi/gostarter/internal/shortly/internal/inbound"
	"github.com/shandysiswandi/gostarter/internal/shortly/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/shortly/internal/service"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/logger"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type Expose struct{}

type Dependency struct {
	RedisDB   *redis.Client
	Config    config.Config
	Validator validation.Validator
	CodecJSON codec.Codec
	Router    *httprouter.Router
	Logger    logger.Logger
}

func New(dep Dependency) (*Expose, error) {
	// redisOutbound := outbound.NewRedisShort(dep.RedisDB, dep.CodecJSON)
	mapOutbound := outbound.NewMapShort()

	getUC := service.NewGet(mapOutbound, dep.Validator, dep.Logger)
	setUC := service.NewSet(mapOutbound, dep.Validator, dep.Logger)

	inbound.RegisterHTTP(dep.Router, &inbound.Endpoint{
		GetUC: getUC,
		SetUC: setUC,
	})

	return &Expose{}, nil
}
