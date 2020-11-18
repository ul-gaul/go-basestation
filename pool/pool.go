package pool

import (
    "github.com/panjf2000/ants/v2"
    log "github.com/sirupsen/logrus"
    "os"
)

// var Backend ants.Pool TODO
var Frontend *ants.Pool

func Release() {
    Frontend.Release()
    ants.Release()
}

func init() {
    var err error
    logger := log.New()
    logger.SetOutput(os.Stderr)
    logger.SetLevel(log.GetLevel())
    
    Frontend, err = ants.NewPool(ants.DefaultAntsPoolSize,
        ants.WithLogger(logger),
        ants.WithPanicHandler(func(err interface{}) {
            log.Panicln(err)
        }),
    )
    if err != nil {
        log.Panicln(err)
    }
}
