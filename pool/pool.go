package pool

import (
    "github.com/panjf2000/ants/v2"
    log "github.com/sirupsen/logrus"
)

// var Backend *ants.Pool // TODO
var Frontend *ants.Pool

func Release() {
    Frontend.Release()
    ants.Release()
}

func init() {
    var err error
    
    Frontend, err = ants.NewPool(ants.DefaultAntsPoolSize,
        ants.WithLogger(log.StandardLogger()),
        ants.WithPanicHandler(func(err interface{}) {
            log.Panicln(err)
        }),
    )
    
    if err != nil {
        log.Panicln(err)
    }
}
