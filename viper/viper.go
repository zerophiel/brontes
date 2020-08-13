package viper

import (
	"sync"

	"github.com/spf13/viper"
)

var (
	once sync.Once
)

func InitViper() {
	once.Do(func() {
		viper.SetConfigType("json")
		viper.AddConfigPath("./data/conf")
		viper.SetConfigName("app.conf")
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
	})
}
