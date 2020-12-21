package cmd

import (
    "gioui.org/app"
    "github.com/panjf2000/ants/v2"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "math"
    "os"
    
    "github.com/ul-gaul/go-basestation/constants"
    "github.com/ul-gaul/go-basestation/data/persistence"
    "github.com/ul-gaul/go-basestation/pool"
    "github.com/ul-gaul/go-basestation/ui"
    "github.com/ul-gaul/go-basestation/utils"
)


func preRun(cmd *cobra.Command, _ []string) error {
    if cmd.Flags().Changed("verbose") {
        verbose, err := cmd.Flags().GetCount("verbose")
        if err != nil { return err }
        verbose += 1
        verbose = int(math.Min(float64(verbose), float64(len(log.AllLevels)-1)))
        verbose = int(math.Max(0, float64(verbose)))
        log.SetLevel(log.Level(verbose))
        log.Debugf("Verbosity changed! (Level: %s)", log.GetLevel())
    }
    return nil
}

func run(cmd *cobra.Command, args []string) error {
    csvFlag := cmd.Flag("load-csv")
    if csvFlag.Changed {
        stat, err := os.Stat(csvFlag.Value.String())
        if err != nil {
            return err
        }
        if !stat.Mode().IsRegular() {
            return constants.ErrNotARegularFile
        }
        csvFile, err = os.Open(csvFlag.Value.String())
        if err != nil { return err }
    }
    return nil
}

func postRun(cmd *cobra.Command, args []string) {
    if csvFile != nil {
        ants.Submit(readCsv)
    }
    
    utils.CheckErr(pool.Frontend.Submit(ui.RunGioui))
    app.Title("Gaul - Base Station")
    app.Main()
}

func readCsv() {
    // FIXME - FONCTION TEMPORAIRE
    
    packets, err := persistence.NewCsvPacketReader(csvFile).ReadAll()
    if err != nil {
        log.Fatalln(err)
    }
    _ = packets
}