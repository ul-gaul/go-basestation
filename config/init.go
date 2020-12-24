package config

import (
    "github.com/mitchellh/go-homedir"
    "github.com/spf13/viper"
    
    "github.com/ul-gaul/go-basestation/utils"
)

func Initialize(cfgFile string) {
    var err error
    
    if cfgFile != "" {
        // Use config file from the flag.
        viper.SetConfigFile(cfgFile)
    } else {
        var home, exe string
        
        // Find home directory.
        home, err = homedir.Dir()
        utils.CheckErr(err)
        
        // Find executable path
        exe, err = utils.GetExecutablePath()
        utils.CheckErr(err)
        
        viper.AddConfigPath(home)
        viper.AddConfigPath(".")
        viper.AddConfigPath(exe)
        viper.SetConfigName("basestation")
    }
    
    viper.AutomaticEnv()
    viper.SetEnvPrefix("BASESTATION")
    
    applyDefaults()
    
    utils.CheckErr(viper.ReadInConfig())
    utils.CheckErr(viper.UnmarshalKey("Comms", &Comms))
    utils.CheckErr(viper.UnmarshalKey("Frontend", &Frontend))
}
