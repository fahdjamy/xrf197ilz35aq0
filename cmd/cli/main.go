package main

import (
	"fmt"
	"os"
	"time"
	"xrf197ilz35aq0"
	xrfErr "xrf197ilz35aq0/internal/error"
)

func main() {
	// connect to the Mongo Database
	//dbConnStr, err := mongoUri(config)
	//if err != nil {
	//	logger.Panic(fmt.Sprintf("appStarted=false :: err%s", err.Error()))
	//	return
	//}
	//databaseName := config.Database.Mongo.DatabaseName
	//mongoClient, err := mongo.NewClient(context.Background(), dbConnStr, databaseName)

	//if err != nil {
	//	internalError := xrfErr.Internal{
	//		Err:     err,
	//		Time:    time.Now(),
	//		Source:  "cmd/cli/main",
	//		Message: "failed to connect to mongo",
	//	}
	//	logger.Panic(fmt.Sprintf("appStarted=false :: err%s", internalError.Error()))
	//	return
	//}

	// create a mongo store
	//mongoStore := mongo.NewStore(logger, mongoClient, databaseName, context.Background())
}

func mongoUri(config xrf197ilz35aq0.Config) (string, error) {
	mongoConfig := config.Database.Mongo
	baseUri := os.Getenv(mongoConfig.Uri)
	if baseUri == "" {
		return "", &xrfErr.Internal{
			Time:    time.Now(),
			Source:  "cmd/cli/main#mongoUri",
			Message: "missing environment variable $" + mongoConfig.Uri,
		}
	}
	return fmt.Sprintf("%s?directConnection=%t&retryWrites=%t&w=%s&appName=%s",
		baseUri,
		mongoConfig.DirectConnection,
		mongoConfig.RetryWrites,
		mongoConfig.Acknowledgment,
		mongoConfig.AppName,
	), nil
}
