package db

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"

	conf "brontes/viper"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func getConnectionString() string {
	conf.InitViper()
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	user := viper.GetString("database.user")
	pass := viper.GetString("database.pass")
	dbname := viper.GetString("database.name")
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		user, pass, host, port, dbname)
}

func GetInstance() *gorm.DB {
	once.Do(func() {
		connectionString := getConnectionString()
		db, err := gorm.Open("postgres", connectionString)
		if err != nil {
			panic(err)
		}
		db.DB().SetMaxIdleConns(viper.GetInt("database.MaxIdleConns"))
		instance = db
	})
	return instance
}
