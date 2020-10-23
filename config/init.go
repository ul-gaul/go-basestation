package config

import (
    "github.com/spf13/viper"
    "os"
    "path/filepath"
    
    "github.com/ul-gaul/go-basestation/utils"
)

func getExecutablePath() string {
    exe, err := os.Executable()
    if err != nil {
        exe, err = filepath.EvalSymlinks(exe)
        utils.CheckErr(err)
    }
    return filepath.Dir(exe)
}

func init() {
    viper.SetConfigName("basestation")
    viper.AddConfigPath(".")
    viper.AddConfigPath("$HOME")
    viper.AddConfigPath(getExecutablePath())
    
    viper.AutomaticEnv()
    viper.SetEnvPrefix("BASESTATION")
    
    applyDefaults()
    
    utils.CheckErr(viper.ReadInConfig())
    utils.CheckErr(viper.UnmarshalKey("comms", &Comms))
}
