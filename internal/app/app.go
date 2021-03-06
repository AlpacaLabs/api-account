package app

import (
	"sync"

	"github.com/AlpacaLabs/api-account/internal/async"

	"github.com/AlpacaLabs/api-account/internal/grpc"

	"github.com/AlpacaLabs/api-account/internal/configuration"
	"github.com/AlpacaLabs/api-account/internal/db"
	"github.com/AlpacaLabs/api-account/internal/http"
	"github.com/AlpacaLabs/api-account/internal/service"
	log "github.com/sirupsen/logrus"
)

type App struct {
	config configuration.Config
}

func NewApp(c configuration.Config) App {
	return App{
		config: c,
	}
}

func (a App) Run() {
	dbConn, err := db.Connect(a.config.SQLConfig)
	if err != nil {
		log.Fatalf("failed to dial database: %v", err)
	}
	dbClient := db.NewClient(dbConn)
	svc := service.NewService(a.config, dbClient)

	var wg sync.WaitGroup

	wg.Add(1)
	httpServer := http.NewServer(a.config, svc)
	go httpServer.Run()

	wg.Add(1)
	grpcServer := grpc.NewServer(a.config, svc)
	go grpcServer.Run()

	wg.Add(1)
	go async.HandleConfirmEmailAddressRequest(a.config, svc)

	wg.Add(1)
	go async.HandleConfirmPhoneNumberRequest(a.config, svc)

	wg.Wait()
}
