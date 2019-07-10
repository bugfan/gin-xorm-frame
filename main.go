package main

import (
	"gin-xorm-frame/api"
	"gin-xorm-frame/influxdb"
	"gin-xorm-frame/models"
	"gin-xorm-frame/setting"
	"os"

	"github.com/gogap/logrus"
)

func main() {
	// init sql db
	_, err := models.SetEngine(&models.Config{
		User:     setting.Get("db_user"),
		Password: setting.Get("db_password"),
		Host:     setting.Get("db_host"),
		Name:     setting.Get("db_name"),
		Log:      setting.Get("db_log"),
	})
	if err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}
	// init mongo
	// err = mgomodels.SetEngine(&mgomodels.Config{
	// 	Host: setting.Get("mongo_url"),
	// })
	// if err != nil {
	// 	logrus.Error(err)
	// 	os.Exit(-2)
	// }

	// init redis
	// _, err = redis.SetEngine(&redis.Config{
	// 	Addr:     setting.Get("redis_addr"),
	// 	Password: setting.Get("redis_password"),
	// 	PoolSize: setting.GetInt("redis_pool_size"),
	// 	DB:       setting.GetInt("redis_index"),
	// })
	// if err != nil {
	// 	logrus.Error(err)
	// 	os.Exit(-3)
	// }

	influxdb.I, err = influxdb.NewClient(
		setting.Get("influx_addr"),
		setting.Get("influx_username"),
		setting.Get("influx_password"),
		setting.Get("influx_db"),
		setting.Get("influx_precision"),
	)
	if err != nil {
		logrus.Error(err)
		os.Exit(-4)
	}

	api, err := api.NewAPIBackend(models.GetEngine())
	if err != nil {
		logrus.Fatal(err)
	}
	if setting.Get("ishttps") == "true" {
		pem := setting.Get("ssl_pem")
		key := setting.Get("ssl_key")
		api.G.RunTLS(setting.Get("server_addr")+":"+setting.Get("server_port"), pem, key)
	} else {
		api.G.Run(setting.Get("server_addr") + ":" + setting.Get("server_port"))
	}
}
