package main

import (
	"fmt"
	"os"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/cmd/cli/dependency"
)

func main() {
	health := xrf197ilz35aq0.NewHealth()
	logger, err := dependency.NewZap(xrf197ilz35aq0.WARN, true, os.Stdout)

	if err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("Current application version '%s'", health.Version()))
}
