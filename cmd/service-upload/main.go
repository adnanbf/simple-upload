package main

import (
	"net/http"
	"os"
	"simple-upload/cmd/service-upload/handler"
	"simple-upload/pkg/config"
	"simple-upload/pkg/util"

	log "github.com/sirupsen/logrus"
)

func main() {
	// config
	log.SetLevel(log.TraceLevel)
	log.SetOutput(os.Stderr)
	if env := os.Getenv("CONFIG_ENV"); env == "" {
		config.LoadEnv("config.local.env")
	}
	config.ShowListEnvs()

	// register handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.NewIndex().GetIndex)
	mux.HandleFunc("/upload", handler.NewUpload().UploadFile)

	serviceHost := os.Getenv("HOST")
	port := ":" + util.GetPortFromHost(serviceHost)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
}
