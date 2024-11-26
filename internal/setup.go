package internal

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"xrf197ilz35aq0/internal/random"
)

const (
	DevelopEnv    = "dev"
	LiveEnv       = "live"
	StagingEnv    = "staging"
	ProductionEnv = "production"
	environment   = "XRF_ENV"
)

type Environment struct {
	Name    string
	LogMode bool
}

func GenerateRequestId() string {
	uniqueStr, err := random.TimeBasedString(time.Now().Unix(), 21)
	if err != nil {
		return strconv.Itoa(int(random.PositiveInt64()))
	}

	uniqueInt64 := random.PositiveInt64()
	uniqueInt64Str := strconv.Itoa(int(uniqueInt64))

	if len(uniqueInt64Str) > 10 {
		uniqueInt64Str = uniqueInt64Str[2:]
	}

	partStr := uniqueStr[0:12]

	return fmt.Sprintf("%s.%s", uniqueInt64Str, partStr)
}

func GetEnvironment() Environment {
	env := os.Getenv(environment)
	if env == "" {
		env = DevelopEnv
	}

	switch env {
	case StagingEnv:
		return Environment{
			Name:    "staging",
			LogMode: true,
		}
	case ProductionEnv, LiveEnv:
		return Environment{
			Name:    "production",
			LogMode: false,
		}
	default:
		return Environment{
			Name:    "dev",
			LogMode: true,
		}
	}
}
