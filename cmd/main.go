package main

import (
	"MScannot206/pkg/login"
	"MScannot206/shared/config"
	"MScannot206/shared/server"
	"context"
	"errors"
	"log"
)

func main() {
	var errs error

	ctx := context.Background()

	webServerCfg := &config.WebServerConfig{
		Port: 8080,

		MongoUri: "mongodb://localhost:27017/",
	}

	web_server, err := server.NewWebServer(
		ctx,
		webServerCfg,
	)

	if err != nil {
		panic(err)
	}

	if login_service, err := login.NewLoginService(
		web_server.GetContext(),
		web_server.GetRouter(),
	); err != nil {
		errs = errors.Join(errs, err)
		log.Println(err)
	} else {
		if err := web_server.AddService(login_service); err != nil {
			errs = errors.Join(errs, err)
			log.Println(err)
		}
	}

	if errs != nil {
		panic(errs)
	}

	if err := web_server.Init(); err != nil {
		panic(err)
	}

	if err := web_server.Start(); err != nil {
		panic(err)
	}
}
