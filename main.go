package main

import (
    "gioui.org/app"
    log "github.com/sirupsen/logrus"
    "os"
    
    _ "github.com/ul-gaul/go-basestation/config"
    "github.com/ul-gaul/go-basestation/pool"
    "github.com/ul-gaul/go-basestation/ui"
)

func init() {
    log.SetOutput(os.Stderr)
    log.SetLevel(log.DebugLevel)
}

func main() {
    defer pool.Release()
    
    if err := pool.Frontend.Submit(ui.RunGioui); err != nil {
        log.Panicln(err)
    }
    app.Title("Gaul - Base Station")
    app.Main()
}