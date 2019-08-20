package newrelic

import (
    "sync"

    "github.com/gin-gonic/gin"
    newrelic "github.com/newrelic/go-agent"
    "github.com/spf13/viper"
)

var (
    instance newrelic.Application
    once     sync.Once
)

func GetInstance() (newrelic.Application,error) {
    once.Do(func() {
        nrAPMName := viper.GetString("newrelic.appname")
        nrLicense := viper.GetString("newrelic.license")
        config := newrelic.NewConfig(nrAPMName, nrLicense)
        Nr, err := newrelic.NewApplication(config)
        if err != nil {
            panic("Error connect to newrelic")
            return
        }
        instance = Nr
    })
    return instance,nil
}

const (
    // NewRelicTxnKey is the key used to retrieve the NewRelic Transaction from the context
    NewRelicTxnKey = "NewRelicTxnKey"
)

// NewRelicMonitoring is a middleware that starts a newrelic transaction, stores it in the context, then calls the next handler
func NewRelicMonitoring(app newrelic.Application) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        txn := app.StartTransaction(ctx.Request.URL.Path, ctx.Writer, ctx.Request)
        defer txn.End()
        ctx.Set(NewRelicTxnKey, txn)
        ctx.Next()
    }
}
