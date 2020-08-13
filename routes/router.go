package routes

import "github.com/gin-gonic/gin"

func CreateNewRouter(app *gin.Engine, relativePath string) *gin.RouterGroup {
	route := (*app).Group(relativePath)
	return route
}
