package main

import (
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
	config, err := conf.GetConfig("pre_live")
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



