package main

import (
	"os"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/log"
)

func main() {
	health := xrf197ilz35aq0.NewHealth()
	logger, err := log.New(os.Stdout, log.DEBUG)

	if err != nil {
		panic(err)
	}

	logger.InfoF("Current application version '%s'", health.Version())
}
