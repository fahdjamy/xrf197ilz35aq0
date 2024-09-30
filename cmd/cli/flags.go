package main

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
)

const (
	maxRetries = 2
	retryAfter = 1 * time.Second
)

var newUserFile = flag.String("New User", "nu", "This stores a new users information")

type Flags struct {
	log        xrf197ilz35aq0.Logger
	otherFlags []string
}

var ParseFlagExtErr = &xrfErr.External{
	Source: "cmd/cli/flags#Parse",
}

var ParseFlagIntErr = &xrfErr.Internal{
	Source: "cmd/cli/flags#Parse",
}

func (f *Flags) Parse() (*exchange.UserRequest, error) {
	flag.Parse()
	now := time.Now()
	ParseFlagExtErr.Time = now
	ParseFlagIntErr.Time = now

	//var newUserReq
	newUserReq := &exchange.UserRequest{}

	if *newUserFile != "" {
		if filepath := flag.Arg(0); filepath != "" {
			// retrieve information about the json file to capture user info from
			fileInfo, err := os.Stat("my_file.txt")
			if err != nil {
				return nil, ParseFlagExtErr.WithErr("can't retrieve the files info", err)
			}
			if fileInfo.IsDir() || !strings.HasSuffix(fileInfo.Name(), ".json") {
				return nil, ParseFlagExtErr.NoErr("not a json file")
			}
			f.log.Info(fmt.Sprintf("event=parseFlag :: flag=%s :: filepath=%s", *newUserFile, filepath))

			// open json file
			userJson, err := os.Open(*newUserFile)
			if err != nil {
				return nil, ParseFlagExtErr.WithErr(fmt.Sprintf("can't open user json file: %v", err), err)
			}
			defer func() {
				if err := xrf197ilz35aq0.CloseFileWithRetry(userJson, maxRetries, retryAfter); err != nil {
					fmt.Println(err)
				}
			}()

			// read our opened jsonFile as a byte array.
			byteValue, err := io.ReadAll(userJson)
			if err != nil {
				return nil, ParseFlagIntErr.WithErr(fmt.Sprintf("can't read user json file: %v", err), err)
			}

			// unmarshal our bytes which contains user info
			if err := json.Unmarshal(byteValue, newUserReq); err != nil {
				return nil, ParseFlagIntErr.WithErr(fmt.Sprintf("can't unmarshal user json file: %v", err), err)
			}
		}
	}

	return newUserReq, nil
}

func NewFlags(logger xrf197ilz35aq0.Logger) Flags {
	return Flags{
		log: logger,
	}
}
