package main

import (
	"github.com/donohutcheon/gowebserver/state"
	"github.com/donohutcheon/gowebserver/state/facotory"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
)

var (
	//CertFile environment variable for CertFile
	CertFile = os.Getenv("CERT_FILE")
	//KeyFile environment variable for KeyFile
	KeyFile = os.Getenv("KEY_FILE")
	//ServiceAddress address to listen on

)

func main() {
	logger := log.New(os.Stdout, "server ", log.LstdFlags|log.Lshortfile)

	// Load environment from .env file for development.
	err := godotenv.Load()
	if err != nil {
		logger.Printf("Could not load environment files. %s", err.Error())
	}

	mode := os.Getenv("ENVIRONMENT")

	mainThreadWG := new(sync.WaitGroup)
	var serverState *state.ServerState
	if mode == "prod" {
		serverState, err = facotory.NewForProduction(logger, mainThreadWG)
	} else {
		serverState, err = facotory.NewForStaging(logger, mainThreadWG)
	}
	if err != nil {
		logger.Printf("failed to initialize production state %s", err.Error())
		return
	}

	mainThreadWG.Wait()
	serverState.Logger.Printf("graceful shutdown")
}