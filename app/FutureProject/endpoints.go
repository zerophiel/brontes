package FutureProject

import (
	"github.com/gin-gonic/gin"
)

func Endpoint(route *gin.RouterGroup) {
	route.POST("/records/", RecordInsert)
	route.GET("/record-count/", CountRecordAll)
}
