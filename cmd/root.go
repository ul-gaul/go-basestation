package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "io"
    "os"
    "path/filepath"
    
    "github.com/ul-gaul/go-basestation/config"
    "github.com/ul-gaul/go-basestation/utils"
)

var (
    cfgFile string
    csvFile io.Reader
)

// rootCmd represents the base command when called without any subcommands
var rootCmd *cobra.Command

func init() {
    cmdName := filepath.Base(os.Args[0])
    
    rootCmd = &cobra.Command{
        Use: cmdName,
        // TODO change descriptions
        Short: "A brief description of your application",
        Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    
        Example: fmt.Sprintf(""+
            "  %s -v \t# Log level: ERROR\n"+
            "  %s -vv \t# Log level: INFO (default)\n"+
            "  %s -vvv \t# Log level: DEBUG\n",
            cmdName, cmdName, cmdName),
            
        PersistentPreRunE: preRun,
    
        RunE: run,
        PersistentPostRun: postRun,
    }
    
    cobra.OnInitialize(initConfig)
    
    // Here you will define your flags and configuration settings.
    // Cobra supports persistent flags, which, if defined here,
    // will be global for your application.
    rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
    rootCmd.PersistentFlags().CountP("verbose", "v", "verbose output / log level")
    
    // Cobra also supports local flags, which will only run
    // when this action is called directly.
    // rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
    rootCmd.Flags().StringP("load-csv", "l", "", "Load data from file")
    utils.CheckErr(rootCmd.MarkFlagFilename("load-csv", "csv"))
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
    // rootCmd.AddCommand(cmds...)
    utils.CheckErr(rootCmd.Execute())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
    config.Initialize(cfgFile)
}
