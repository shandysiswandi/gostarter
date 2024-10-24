package gallery

import (
	"database/sql"

	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
	"github.com/shandysiswandi/gostarter/internal/gallery/internal/inbound"
	"github.com/shandysiswandi/gostarter/internal/gallery/internal/outbound"
	"github.com/shandysiswandi/gostarter/internal/gallery/internal/service"
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
	Validator validation.Validator
	CodecJSON codec.Codec
	Router    *httprouter.Router
	Telemetry *telemetry.Telemetry
}

func New(dep Dependency) (*Expose, error) {
	sql := outbound.NewSQLGallery(dep.Database)

	uploadUC := service.NewUpload(sql)
	getImageUC := service.NewGetImage(sql)

	inbound.RegisterHTTP(dep.Router, &inbound.Endpoint{
		UploadUC:   uploadUC,
		GetImageUC: getImageUC,
	})

	return &Expose{}, nil
}
