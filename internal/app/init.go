package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"github.com/shandysiswandi/gostarter/pkg/clock"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

func (a *App) initSTDLog() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func (a *App) initConfig() {
	cfg, err := config.NewViperConfig("config/config.yaml")
	if err != nil {
		log.Fatalln("failed to init config", err)
	}

	a.config = cfg
}

func (a *App) initDatabase() {
	maxOpen := a.config.GetInt(`database.max.open`)
	maxIdle := a.config.GetInt(`database.max.idle`)
	maxLifetime := a.config.GetInt(`database.max.lifetime`)
	maxIdletime := a.config.GetInt(`database.max.idletime`)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		a.config.GetString(`database.user`),
		a.config.GetString(`database.pass`),
		a.config.GetString(`database.host`),
		a.config.GetString(`database.port`),
		a.config.GetString(`database.name`),
	)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	val.Encode()

	database, err := sql.Open("mysql", fmt.Sprintf("%s?%s", dsn, val.Encode()))
	if err != nil {
		log.Fatalln("failed to open database", err)
	}

	if err := database.Ping(); err != nil {
		log.Fatalln("failed to ping database connection", err)
	}

	database.SetMaxOpenConns(int(maxOpen))
	database.SetMaxIdleConns(int(maxIdle))
	database.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Minute)
	database.SetConnMaxIdleTime(time.Duration(maxIdletime) * time.Minute)

	a.database = database
}

func (a *App) initRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr: a.config.GetString("redis.addr"),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalln("failed to init redis", err)
	}

	a.redisdb = rdb
}

func (a *App) initHTTPRouter() {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("content-type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(map[string]string{"code": "40400", "message": "Endpoint not found"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("content-type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		err := json.NewEncoder(w).Encode(map[string]string{"code": "40500", "message": "Method endpoint not allowed"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	a.httpRouter = router
}

func (a *App) initHTTPServer() {
	handler := cors.Default().Handler(a.httpRouter)

	a.httpServer = &http.Server{
		Addr:              a.config.GetString("server.address.http"),
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}

func (a *App) initLibraries() {
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

func (a *App) initTasks() {
	for _, runnable := range a.runables {
		if err := runnable.Start(); err != nil {
			log.Fatalln("failed to init initTasks", err)
		}
	}
}
