package main

import (
    log "github.com/sirupsen/logrus"
    "os"
    
    "github.com/ul-gaul/go-basestation/cmd"
    _ "github.com/ul-gaul/go-basestation/config"
    "github.com/ul-gaul/go-basestation/constants"
    "github.com/ul-gaul/go-basestation/pool"
)

func init() {
    log.SetOutput(os.Stderr)
    log.SetLevel(constants.DefaultLogLevel)
}

func main() {
    defer pool.Release()
    cmd.Execute()
    // See cmd/run.go
}
