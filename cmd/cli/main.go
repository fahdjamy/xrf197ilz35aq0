package main

import (
	"fmt"
	"os"
	"time"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/cmd/cli/dependency"
	"xrf197ilz35aq0/random"
)

func main() {
	health := xrf197ilz35aq0.NewHealth()
	logger, err := dependency.NewZap(xrf197ilz35aq0.WARN, true, os.Stdout)

	if err != nil {
		panic(err)
	}
	processId, err := generateRequestId()
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info(fmt.Sprintf("requestId=%s, application version '%s'", processId, health.Version()))
}

func generateRequestId() (string, error) {
	uniqueStr, err := random.TimeBasedString(time.Now().Unix(), 21)
	if err != nil {
		return "", err
	}

	uniqueInt64 := random.PositiveInt64()
	partStr := uniqueStr[0:7]

	return fmt.Sprintf("%s.%d", partStr, uniqueInt64), nil
}
