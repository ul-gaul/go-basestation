package utils

import (
    log "github.com/sirupsen/logrus"
    "os"
    "path/filepath"
)

// CheckErr log l'erreur et panic si err n'est pas nil.
func CheckErr(err error) {
    if err != nil {
        log.Panic(err)
    }
}


// GetExecutablePath
func GetExecutablePath() (string, error) {
    exe, err := os.Executable()
    if err != nil {
        exe, err = filepath.EvalSymlinks(exe)
        if err != nil {
            return exe, err
        }
    }
    return filepath.Dir(exe), nil
}