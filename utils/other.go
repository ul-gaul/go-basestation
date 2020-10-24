package utils

import (
    "log"
)

// CheckErr log l'erreur et panic si err n'est pas nil.
func CheckErr(err error) {
    if err != nil {
        log.Panicln(err)
    }
}