package config

import (
    "github.com/spf13/viper"
    "log"
    "os"
    "path/filepath"
)

func checkErr(err error) {
    if err != nil {
        log.Panicln(err)
    }
}

func getExecutablePath() string {
    exe, err := os.Executable()
    if err != nil {
        exe, err = filepath.EvalSymlinks(exe)
        checkErr(err)
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
    
    checkErr(viper.ReadInConfig())
    checkErr(viper.UnmarshalKey("comms", &Comms))
}
