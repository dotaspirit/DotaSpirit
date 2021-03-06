package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"gopkg.in/gographics/imagick.v3/imagick"
)

var (
	appconfig = appConfig{}
	dao       = dbDao{}
)

func init() {
	loadConfig("config/app.json", &appconfig)
	dao.Server = appconfig.DBServer
	dao.Database = appconfig.DBName
	dao.connect()
	imagick.Initialize()
	defer imagick.Terminate()
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", webhoookHandler)
	srv := &http.Server{
		Addr:    ":3682",
		Handler: cors.Default().Handler(mux),
	}
	log.Fatal(srv.ListenAndServe())
}
