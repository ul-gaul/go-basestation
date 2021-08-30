package main

import (
    log "github.com/sirupsen/logrus"
    "github.com/ul-gaul/go-basestation/cmd"
    "github.com/ul-gaul/go-basestation/controller"
    "os"
    
    _ "github.com/ul-gaul/go-basestation/cfg"
    "github.com/ul-gaul/go-basestation/pool"
)

func init() {
    log.SetOutput(os.Stderr)
    log.SetLevel(log.InfoLevel)
}

func main() {
    defer controller.Shutdown()
    defer pool.Release()
    cmd.Execute()
}
