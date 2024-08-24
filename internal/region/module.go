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

/*
- /regions/search
- /regions/search?ids=&by=
- /regions/search?ids=
- /regions/search?by=
# These paths will return a 200 status with an empty list.

- /regions/search?by=provinces
- /regions/search?by=cities
- /regions/search?by=districts
- /regions/search?by=villages
# These paths will return a 200 status with all records from the specified table.

- /regions/search?ids=1
- /regions/search?ids=1,2
- /regions/search?ids=1,2,......n
# These paths will return a 200 status with all records corresponding to the IDs,
searching across all tables from provinces to villages.

- /regions/search?ids=1,2&by=provinces
- /regions/search?ids=1,2&by=cities
- /regions/search?ids=1,2&by=districts
- /regions/search?ids=1,2&by=villages
# These paths will return a 200 status with records based on the IDs, searching within
the table specified by the by parameter.

# If the by parameter is not one of provinces, cities, districts, or villages,
the response will be a 400 status.
# If the ids parameter contains non-numeric values, the response will be a 400 status.
*/

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
