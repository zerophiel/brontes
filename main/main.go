package main

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/googollee/go-socket.io"
	_ "github.com/lib/pq"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"brontes/app"
	"brontes/db"
	"brontes/middleware"
	"brontes/newrelic"
	"brontes/nrgorm"
	conf "brontes/viper"
)

var d = db.GetInstance()

func WrapDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		d = nrgorm.Wrap(nrgin.Transaction(c), d)
	}
}

var (
	srv       *http.Server
	socketSrv *socketio.Server
)

func main() {
	conf.InitViper()
	appName := viper.GetString("appname")
	application := gin.Default()
	var fileName = "./data/logs/app.log"
	var aLogFileName = "./data/logs/access.log"
	f, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	a, _ := os.OpenFile(aLogFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	Formatter := new(logrus.JSONFormatter)
	logrus.SetFormatter(Formatter)
	mw := io.MultiWriter(os.Stdout, f)
	logrus.SetOutput(mw)

	aLog := logrus.New()
	aFormatter := new(logrus.JSONFormatter)
	aLog.SetFormatter(aFormatter)
	aMw := io.MultiWriter(a)
	aLog.SetOutput(aMw)

	application.Use(middleware.Logger(aLog))
	if viper.GetBool("newrelic.enabled") {
		app, err := newrelic.GetInstance()
		if err != nil {
			logrus.Error("failed to make new_relic app: %v", err)
		} else {
			application.Use(newrelic.NewRelicMonitoring(app))
			application.Use(WrapDB())
		}
	}
	if viper.GetBool("cors") {
		application.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"*"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"*"},
			AllowCredentials: true,
			AllowWebSockets:  true,

			MaxAge: 12 * time.Hour,
		}))
	}
	app.GenerateAppEndpoint(application, "/"+appName)
	portNum := strconv.Itoa(viper.GetInt("port"))

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				type (
					Result struct {
						Name  string
						Count int
					}
					Request struct {
						Data []Result `json:"data"`
					}
				)
				var (
					result []Result
				)

				d := db.GetInstance()
				stat := d.Raw(`SELECT COUNT (f.name) as Count,f.name from future_programs f group by f.name`).Scan(&result)
				logrus.Info(stat.RowsAffected)
				client := resty.New()
				_, _ = client.R().
					SetHeader("Content-Type", "application/json").
					SetBody(Request{Data: result}).
					Post("'localhost:5556/update-scoreboard")
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	err := application.Run(":" + portNum)

	if err != nil {
		panic("Error to run server")
	}
}
