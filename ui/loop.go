package ui

import (
    "fmt"
    "gioui.org/app"
    "gioui.org/font/gofont"
    "gioui.org/unit"
    "gioui.org/widget/material"
    log "github.com/sirupsen/logrus"
    "gonum.org/v1/plot/plotter"
    "time"
    
    "github.com/ul-gaul/go-basestation/data/packet"
    "github.com/ul-gaul/go-basestation/ui/plotting"
    "github.com/ul-gaul/go-basestation/ui/plotting/lines"
    "github.com/ul-gaul/go-basestation/ui/plotting/ticker"
    "github.com/ul-gaul/go-basestation/ui/views"
    "github.com/ul-gaul/go-basestation/ui/widgets"
    "github.com/ul-gaul/go-basestation/utils"
    
    "gioui.org/io/system"
    "gioui.org/layout"
    "gioui.org/op"
)

var (
    window *app.Window
    theme  *material.Theme
)

func loop() {
    defer log.Exit(0)
    
    // Initilialise la fenêtre de l'application
    window = app.NewWindow(app.Title("GAUL - Base Station"))
    theme = material.NewTheme(gofont.Collection())
    defer window.Close()
    
    // Initialise l'onglet "General"
    generalTab, err := views.NewGeneralTab()
    utils.CheckErr(err)
    
    /************************** **************************/
    drawer, err := plotting.NewPlotDrawer()
    utils.CheckErr(err)
    drawer.Chart().Add(lines.NewOriginLines())
    drawer.Chart().Add(plotter.NewGrid())
    drawer.Chart().X.Tick.Marker = ticker.NewTicker(10, ticker.ContainData)
    
    // Crée un Plotter pour l'onglet "Test"
    plter, err := plotting.NewPlotter(
        // plotting.WithStyleIdx(0),
        // plotting.WithPointStyle(draw.GlyphStyle{Radius: 1.5, Shape: draw.CircleGlyph{}}),
        plotting.WithLegend(fmt.Sprintf("Line #%d", 1)),
        plotting.WithDataLimit(150),
        plotting.WithData(squarePlot(0)...),
    )
    utils.CheckErr(err)
    utils.CheckErr(drawer.AddPlotter(plter))
    
    // Crée la liste d'onglets
    tabs := []widgets.TabChild{
        widgets.Tabbed("General", generalTab.Layout),
        widgets.Tabbed("Motor", generalTab.Layout), // TODO
        widgets.Tabbed("Test", func(gtx layout.Context) layout.Dimensions {
            return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                return layout.UniformInset(unit.Px(10)).Layout(gtx, drawer.Layout)
            })
        }),
    }
    
    /************************** **************************/
    ops := new(op.Ops)
    tabBar := widgets.Tab(theme)
    
    tick := time.NewTicker(100 * time.Millisecond)
    defer tick.Stop()
    
    chData := make(chan packet.PacketList)
    Collector.AddCallback(func(packets packet.PacketList) { chData <- packets })
    generalTab.Plotters()[views.PltAltitude].AppendAll(Collector.Packets().AltitudeData())
    
    for {
        select {
        
        case event := <-window.Events(): // À chaque évènement de l'application
            // Si le type de l'évènement est...
            switch e := event.(type) {
            case system.DestroyEvent: // Fermeture de l'application
                utils.CheckErr(e.Err)
                return
            case system.FrameEvent: // Nouvelle frame/image générée
                // Crée le contenu de la prochaine frame
                gtx := layout.NewContext(ops, e)
                tabBar.Layout(gtx, tabs...)
                
                // Génère la frame (relance une system.FrameEvent)
                e.Frame(gtx.Ops)
            }
            
        case packets := <-chData: // À chaque donnée reçue...
            // ajoute les données au graphique
            generalTab.Plotters()[views.PltAltitude].AppendAll(packets.AltitudeData())
        
        // Ce case permet de simuler des données reçues en temps réel
        case <-tick.C: // À chaque 100 ms...
            // Ajoute 1 point random à plter
            addRandomPoints(plter, 1)
            // Log dans la console le nombre de points que plter contient
            log.Infof("Points: %d", plter.Data().Len())
        }
    }
}
