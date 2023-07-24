package main

import (
	"github.com/ahmetson/sample-service/handler"
	"github.com/ahmetson/service-lib/configuration"
	"github.com/ahmetson/service-lib/independent"
	"github.com/ahmetson/service-lib/log"
)

const webProxyUrl = "github.com/ahmetson/web-proxy"

// sample extension that utilizes calculator
const calcUrl = "github.com/ahmetson/sample-extension"

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

	replier, _ := handler.NewReplier(logger, calcUrl)
	service.AddController("replier", replier)

	service.RequireProxy(webProxyUrl, configuration.DefaultContext)

	err = service.Pipe(webProxyUrl, "replier")
	if err != nil {
		logger.Fatal("service.Pipe", "error", err)
	}

	err = service.Prepare(configuration.IndependentType)
	if err != nil {
		logger.Fatal("service.Prepare", "error", err)
	}

	service.Run()
}
