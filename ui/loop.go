package ui

import (
    "gioui.org/app"
    "github.com/panjf2000/ants/v2"
    
    "github.com/ul-gaul/go-basestation/constants"
    . "github.com/ul-gaul/go-basestation/data/collector"
    "github.com/ul-gaul/go-basestation/data/packet"
    "github.com/ul-gaul/go-basestation/data/parsing"
    . "github.com/ul-gaul/go-basestation/data/persistence"
    "github.com/ul-gaul/go-basestation/pool"
    "github.com/ul-gaul/go-basestation/utils"
)

var (
    running   bool
    collector IDataCollector
    parser    parsing.ISerialPacketParser
    conn      ISerialPacketCommunicator
)

func init() {
    var err error
    
    collector, err = New()
    utils.CheckErr(err)
    
    parser, err = parsing.UnmarshalParser()
    utils.CheckErr(err)
}

func Run() {
    if running {
        return
    }
    running = true
    utils.CheckErr(pool.Frontend.Submit(loop))
    app.Main()
}

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
                collector.AddPackets(pkt)
            case err = <-conn.ErrorChannel():
                // TODO Rendre moins agressif (ex: afficher l'erreur dans l'interface et fermer la connexion)
                utils.CheckErr(err)
            }
        }
    })
}
