package app

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"brontes/app/FutureProject"
	"brontes/db"
	"brontes/models"
	"brontes/routes"
)

var d = db.GetInstance()

func initDB() {
	d.AutoMigrate(models.FutureProgram{})
}

func GenerateAppEndpoint(app *gin.Engine, appName string) {
	initDB()
	d.LogMode(viper.GetBool("database.debug"))
	FutureProject.Endpoint(routes.CreateNewRouter(app, appName+"/future-project"))
}
