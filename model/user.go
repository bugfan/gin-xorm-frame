package model

import "time"

func init() {
	Register(new(User))
}

type User struct {
	ID       int64
	Username string
	Password string
	Created  time.Time
	Updated  time.Time
}
