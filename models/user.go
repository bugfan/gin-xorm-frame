package models

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

func (u *User) Explain() string {
	if u.ID > 3 {
		return "add_extend_column > 3"
	}
	return "add_extend_column"
}
