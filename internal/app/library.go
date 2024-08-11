package app

import (
	"log"

	"github.com/shandysiswandi/gostarter/pkg/clock"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

func (a *App) initlibraries() {
	snow, err := uid.NewSnowflakeNumber()
	if err != nil {
		log.Fatalln("failed to init uid number snowflake", err)
	}

	a.uidnumber = snow
	a.clock = clock.NewClock()
	a.uuid = uid.NewUUIDString()
	a.codecJSON = codec.NewJSONCodec()
	a.codecMsgPack = codec.NewMsgpackCodec()
	a.validator = validation.NewV10Validator()
}
