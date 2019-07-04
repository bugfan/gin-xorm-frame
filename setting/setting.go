package setting

import (
	"gin-xorm-frame/models"
	"log"
	"os"
	"strings"
	"time"
)

var defaults map[string]string

func init() {
	models.Register(new(Setting))
	defaults = map[string]string{
		"db_user":     "root",
		"db_password": "",
		"db_host":     "127.0.0.1:3306",
		"db_name":     "scaffold",
		"db_log":      "xorm.log",
		"server_port": "9997",
		"server_addr": "0.0.0.0",
		"ssl_key":     "./ssl.key",
		"ssl_pem":     "./ssl.pem",
		"ishttps":     "false",
		"jwt_secret":  "", // "" is use random string
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
