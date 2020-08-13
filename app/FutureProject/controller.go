package FutureProject

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"brontes/db"
	"brontes/models"
)

func RecordInsert(c *gin.Context) {
	var (
		payload models.FutureProgram
	)
	d := db.GetInstance()
	err := c.BindJSON(&payload)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "success": false, "message": err})
		return
	}
	stat := d.Save(&payload)
	logrus.Info(stat.RowsAffected)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "success": true, "message": "success", "data": stat})
}

func CountRecordAll(c *gin.Context) {
	type Result struct {
		Name  string
		Count int
	}
	var (
		result []Result
	)

	d := db.GetInstance()
	stat := d.Raw(`SELECT COUNT (f.name) as Count,f.name from future_programs f group by f.name`).Scan(&result)
	logrus.Info(stat.RowsAffected)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "success": true, "message": "success", "data": result})
}
