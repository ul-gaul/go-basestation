package cmd

import (
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "math"
    "os"
    
    "github.com/ul-gaul/go-basestation/constants"
    "github.com/ul-gaul/go-basestation/controller"
    "github.com/ul-gaul/go-basestation/data/persistence"
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

func run(cmd *cobra.Command, _ []string) error {
    csvFlag := cmd.Flag("load-csv")
    if csvFlag.Changed {
        stat, err := os.Stat(csvFlag.Value.String())
        if err != nil {
            return err
        }
        if !stat.Mode().IsRegular() {
            return constants.ErrNotARegularFile
        }
        
        packets, err := persistence.ReadCsv(csvFlag.Value.String())
        if err != nil { return err }
        controller.Collector().AddPackets(packets...)
    }
    return nil
}