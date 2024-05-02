package main

import (
	"time"

	route "github.com/jordanlanch/modak-test/api/route"
	"github.com/jordanlanch/modak-test/bootstrap"
)

func main() {

	app := bootstrap.App(".env")

	env := app.Env

	timeout := time.Duration(env.ContextTimeout) * time.Second

	router := route.Setup(env, timeout, app.Rdb)

	router.Run(env.ServerAddress)
}
