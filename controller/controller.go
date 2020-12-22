package controller

import (
    "gioui.org/app"
    "gioui.org/font/gofont"
    "gioui.org/widget/material"
    "github.com/panjf2000/ants/v2"
    log "github.com/sirupsen/logrus"
    
    "github.com/ul-gaul/go-basestation/constants"
    "github.com/ul-gaul/go-basestation/data/collector"
    "github.com/ul-gaul/go-basestation/data/packet"
    "github.com/ul-gaul/go-basestation/data/parsing"
    . "github.com/ul-gaul/go-basestation/data/persistence"
    "github.com/ul-gaul/go-basestation/pool"
    "github.com/ul-gaul/go-basestation/ui"
    "github.com/ul-gaul/go-basestation/utils"
)

var (
    running bool
    data    collector.IDataCollector
    window  *app.Window
    theme   *material.Theme
    parser  parsing.ISerialPacketParser
    conn    ISerialPacketCommunicator
)

func Window() *app.Window                   { return window }
func Theme() *material.Theme                { return theme }
func Collector() collector.IDataCollector   { return data }
func Connection() ISerialPacketCommunicator { return conn }

func Connect(port string) error {
    // TODO Ajouter la possibilité de fermer une connexion
    // TODO Permettre de créer une autre connexion uniquement si la précédente est fermée
    if conn != nil {
        return constants.ErrConnectionAlreadyOpen
    }
    
    conn = NewSerialPacketCommunicator(port, parser)
    if err := conn.Start(); err != nil {
        return err
    }
    
    return ants.Submit(func() {
        var pkt packet.RocketPacket
        var err error
        for {
            select {
            case pkt = <-conn.RocketPacketChannel():
                Collector().AddPackets(pkt)
            case err = <-conn.ErrorChannel():
                // TODO Rendre moins agressif (ex: afficher l'erreur dans l'interface et fermer la connexion)
                utils.CheckErr(err)
            }
        }
    })
}

func init() {
    var err error
    
    data, err = collector.New()
    utils.CheckErr(err)
    
    parser, err = parsing.UnmarshalParser()
    utils.CheckErr(err)
}

func Run() {
    if running {
        return
    }
    running = true
    
    utils.CheckErr(pool.Frontend.Submit(func() {
        defer log.Exit(0)
        
        window = app.NewWindow(app.Title("GAUL - Base Station"))
        theme = material.NewTheme(gofont.Collection())
        defer window.Close()
        
        utils.CheckErr(ui.Loop(window, theme))
    }))
    
    app.Main()
}
