package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ul-gaul/go-basestation/controller"
	"github.com/ul-gaul/go-basestation/ui/engine"
	"github.com/ul-gaul/go-basestation/ui/vars"
	"math"
)

func preRun(cmd *cobra.Command, _ []string) error {
	if cmd.Flags().Changed("verbose") {
		verbose, err := cmd.Flags().GetCount("verbose")
		if err != nil {
			return err
		}
		verbose += 1
		verbose = int(math.Min(float64(verbose), float64(len(log.AllLevels)-1)))
		verbose = int(math.Max(0, float64(verbose)))
		log.SetLevel(log.Level(verbose))
		log.Debugf("Verbosity changed! (Level: %s)\n", log.GetLevel())
	}
	return nil
}

func run(cmd *cobra.Command, _ []string) error {
	csvFlag := cmd.Flag("load-csv")
	if csvFlag.Changed {
		return controller.LoadFile(csvFlag.Value.String())
	}
	return nil
}

func postRun(_ *cobra.Command, _ []string) {
	vars.Initialize()
	engine.Run()
}
