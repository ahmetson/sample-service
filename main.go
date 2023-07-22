package main

import (
	"github.com/ahmetson/sample-service/handler"
	"github.com/ahmetson/service-lib/configuration"
	"github.com/ahmetson/service-lib/independent"
	"github.com/ahmetson/service-lib/log"
)

const webProxyUrl = "github.com/ahmetson/web-proxy"

func main() {
	logger, err := log.New("sample", false)
	if err != nil {
		log.Fatal("log.New", "error", err)
	}
	appConfig, err := configuration.New(logger)
	if err != nil {
		logger.Fatal("configuration.New", "error", err)
	}

	service, err := independent.New(appConfig, logger.Child("service"))
	if err != nil {
		logger.Fatal("independent.New", "error", err)
	}

	replier, _ := handler.NewReplier(logger)
	service.AddController("replier", replier)

	service.RequireProxy(webProxyUrl)

	err = service.Pipe(webProxyUrl, "replier")
	if err != nil {
		logger.Fatal("service.Pipe", "error", err)
	}

	err = service.Prepare()
	if err != nil {
		logger.Fatal("service.Prepare", "error", err)
	}

	service.Run()
}
