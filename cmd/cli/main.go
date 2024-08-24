package main

import (
	"fmt"
	"os"
	"strconv"
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
	uniqueInt64Str := strconv.Itoa(int(uniqueInt64))

	if len(uniqueInt64Str) > 10 {
		uniqueInt64Str = uniqueInt64Str[:11]
	}

	partStr := uniqueStr[0:12]

	return fmt.Sprintf("%s.%s", partStr, uniqueInt64Str), nil
}
