package main

import (
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/shandysiswandi/gostarter/internal/app"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	application := app.New()
	wait := application.Start()
	<-wait
	application.Stop(ctx)
}
