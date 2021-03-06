package setting

import (
	"gin-xorm-frame/models"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var defaults map[string]string

func init() {
	models.Register(new(Setting))
	defaults = map[string]string{
		"db_user":          "root",
		"db_password":      "",
		"db_host":          "127.0.0.1:3306",
		"db_name":          "scaffold",
		"db_log":           "xorm.log",
		"mongo_url":        "127.0.0.1",
		"redis_addr":       "127.0.0.1",
		"redis_password":   "",
		"redis_pool_size":  "100",
		"redis_index":      "0",
		"influx_addr":      "http://127.0.0.1:8086",
		"influx_username":  "zhaoxy",
		"influx_password":  "zhaoxy0219",
		"influx_db":        "scaffold",
		"influx_precision": "ms",
		"server_port":      "9997",
		"server_addr":      "0.0.0.0",
		"ssl_key":          "./ssl.key",
		"ssl_pem":          "./ssl.pem",
		"ishttps":          "false",
		"jwt_secret":       "", // "" is use random string
	}
}

type Setting struct {
	Key     string `xorm:"pk"`
	Value   string
	Created time.Time `xorm:"CREATED"`
	Updated time.Time `xorm:"UPDATED"`
	Deleted time.Time `xorm:"deleted"`
}

func getDefault(key string) string {
	return defaults[key]
}

func Get(key string) string {
	// 1.get from env
	env := strings.TrimSpace(os.Getenv(strings.ToUpper(key)))
	if env != "" {
		return env
	}
	// 2.then from db
	x := models.GetEngine()
	if x == nil {
		return getDefault(key)
	}
	s := new(Setting)
	has, err := x.ID(key).Get(s)
	if err != nil || !has {
		// 3.finally from defaults
		return getDefault(key)
	}
	return s.Value
}
func GetInt(key string) int {
	v, _ := strconv.Atoi(Get(key))
	return v
}
func Set(key, value string) {
	x := models.GetEngine()
	if x == nil {
		return
	}
	s := new(Setting)
	has, err := x.ID(key).Get(s)
	if err != nil {
		log.Printf("setting %s fail %s", key, err.Error())
	}
	if !has {
		s.Key = key
		s.Value = value
		x.Insert(s)
	} else {
		s.Value = value
		x.ID(key).Cols("value").Update(s)
	}
}
