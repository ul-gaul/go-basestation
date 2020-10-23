package main

import (
    "gioui.org/app"
    "github.com/panjf2000/ants/v2"
    _ "github.com/ul-gaul/go-basestation/config"
    "github.com/ul-gaul/go-basestation/ui"
    "log"
)

func main() {
    defer ants.Release()
    
    const useFyne = false
    
    if useFyne {
        ui.RunFyne()
    } else {
        if err := ants.Submit(ui.RunGioui); err != nil {
            log.Panicln(err)
        }
        app.Main()
    }
}