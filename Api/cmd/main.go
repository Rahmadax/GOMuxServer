package main

import (
	"flag"
	"fmt"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg"
	"net/http"
	"os"
)

// Setup
func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)

		os.Exit(1)
	}
}

func run() error {
	configFilePath := flag.String("config", "Api/config/pre_live.yml", "Path to config file")
	config, err := conf.GetConfig(*configFilePath)
	if err != nil {
		return err
	}

	app, err := pkg.NewApp(*config)
	if err != nil {
		return err
	}

	println("Listening on port - " + config.Server.HostPort)
	return http.ListenAndServe(":"+config.Server.HostPort, app.Router)
}



