package main

import (
	"gin-xorm-frame/api"
	"gin-xorm-frame/models"
	"gin-xorm-frame/setting"
	"os"

	"github.com/gogap/logrus"
)

func main() {
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
