package main

import (
	"fmt"
	"log"

	"net/http"

	"srds.com/srdsapi/data"
	"srds.com/srdsapi/internal/driver"
	"srds.com/srdsapi/logger"
)

const webPort = "8088"

type Config struct {
	DSN  string
	Repo data.Repository
}

func main() {

	app := Config{
		DSN: "host=host.docker.internal port=5432 user=postgres password=password dbname=marikadb sslmode=disable timezone=UTC connect_timeout=5",
	}

	db, err := driver.ConnectSQL(app.DSN)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	defer db.SQL.Close()
	logger.Info("Connected to Postgres DB")

	// http.FS can be used to create a http Filesystem
	var assets = http.FS(assetsFS)
	fs := http.FileServer(assets)

	// Serve static files
	http.Handle("/assets/", fs)

	logger.Info(fmt.Sprintf("Starting front end service on port %s ", webPort))

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		log.Panic(err)
	}

}
