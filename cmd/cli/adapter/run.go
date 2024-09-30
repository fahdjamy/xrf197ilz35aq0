package adapter

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"xrf197ilz35aq0"
	"xrf197ilz35aq0/core/exchange"
	xrfErr "xrf197ilz35aq0/internal/error"
	"xrf197ilz35aq0/storage"
)

const (
	maxRetries = 2
	retryAfter = 1 * time.Second
)

var newUserFile = flag.String("nu", "", "This stores a new users information")

type App struct {
	log   xrf197ilz35aq0.Logger
	mongo storage.Store
}

var parseFlagExtErr = &xrfErr.External{
	Source: "cmd/cli/run#start",
}

var parseFlagIntErr = &xrfErr.Internal{
	Source: "cmd/cli/run#start",
}

func (app *App) Start() {
	flag.Parse()
	now := time.Now()
	parseFlagExtErr.Time = now
	parseFlagIntErr.Time = now

	userAdapter := NewUserAdapter(app.mongo, app.log)

	newUserReq := &exchange.UserRequest{}

	app.log.Info("Provide flags to be evaluated")
	if *newUserFile != "" {
		app.log.Info(fmt.Sprintf("event=newUser :: filePath=%s", *newUserFile))

		// retrieve information about the json file to capture user info from
		fileInfo, err := os.Stat(*newUserFile)
		if err != nil {
			app.log.Error(fmt.Sprintf("err=can't retrieve the files info"))
			return
		}
		if fileInfo.IsDir() || !strings.HasSuffix(fileInfo.Name(), ".json") {
			app.log.Error(fmt.Sprintf("err=file %s is not a json file", newUserFile))
			return
		}
		app.log.Info(fmt.Sprintf("event=parseFlag :: flag=%s :: filepath=%s", *newUserFile, newUserFile))

		// open json file
		userJson, err := os.Open(*newUserFile)
		if err != nil {
			app.log.Error(fmt.Sprintf("err=can't open json file %s", *newUserFile))
			return
		}
		defer func() {
			if err := xrf197ilz35aq0.CloseFileWithRetry(userJson, maxRetries, retryAfter); err != nil {
				app.log.Error(fmt.Sprintf("err=can't close file %s", *newUserFile))
			}
		}()

		// read our opened jsonFile as a byte array.
		byteValue, err := io.ReadAll(userJson)
		if err != nil {
			app.log.Error(fmt.Sprintf("err=can't read json file %s", *newUserFile))
			return
		}

		// unmarshal our bytes which contains user info
		if err := json.Unmarshal(byteValue, newUserReq); err != nil {
			app.log.Error(fmt.Sprintf("err=can't unmarshal json file %s", *newUserFile))
			return
		}

		// create user
		userAdapter.CreateUser(newUserReq)
	}
}

//func logErr(err error, log xrf197ilz35aq0.Logger) {
//	if errors.Is(err, parseFlagExtErr) {
//		log.Panic(err.Error())
//		return
//	}
//	if errors.Is(err, parseFlagIntErr) {
//		log.Error(err.Error())
//		return
//	}
//}

func NewApp(logger xrf197ilz35aq0.Logger, store storage.Store) App {
	return App{
		log:   logger,
		mongo: store,
	}
}
