package mgomodels

import (
	"time"
)

type accessLog struct {
	Level, Url, Data string
	Time             time.Time
	Tags             map[string]string
}

func (s *accessLog) Save() error {
	session := session.Copy() // or GetSession()
	defer session.Close()
	return session.DB(MONGODB_NAME).C("access_log").Insert(s)
}

func NewAccessLog(Level, Url, Data string, Tags map[string]string) *accessLog {
	return &accessLog{
		Level: Level,
		Url:   Url,
		Data:  Data,
		Time:  time.Now(),
		Tags:  Tags,
	}
}
