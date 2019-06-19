package main

import (
	"os"
	"scaffold/api"
	"scaffold/model"
	"scaffold/setting"

	"github.com/gogap/logrus"
)

func main() {
	_, err := model.SetEngine(&model.Config{
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

	api, err := api.NewAPIBackend(model.GetEngine())
	if err != nil {
		logrus.Fatal(err)
	}
	api.G.Run(":9999")
}
