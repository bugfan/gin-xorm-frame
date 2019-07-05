package mgomodels

import mgo "gopkg.in/mgo.v2"

var (
	session *mgo.Session
)

const (
	MONGODB_NAME = "scaffold"
)

type Config struct {
	Host string
}

func GetSession() *mgo.Session {
	return session.Copy()
}

func SetEngine(config *Config) error {
	var err error
	session, err = mgo.Dial(config.Host)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Eventual, false)
	return nil
}
