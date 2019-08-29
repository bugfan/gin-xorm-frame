package models

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"time"

	"log"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

func init() {
	tables = append(tables, new(Admin))
}

type Admin struct {
	ID                int64
	Username          string `xorm:"index"`
	EncryptedPassword string `json:"-"`
	Salt              string `json:"-"`
	Password          string `xorm:"-" json:"-"`

	Created time.Time `xorm:"CREATED"`
	Updated time.Time `xorm:"UPDATED"`
	Deleted time.Time `xorm:"deleted"`
}

func FindAdmin(username string) *Admin {
	u := &Admin{
		Username: username,
	}
	has, _ := x.Get(u)
	if has {
		return u
	}
	return nil
}

func NewAdmin(username, password string) *Admin {
	u := &Admin{
		Username: username,
	}
	u.SetPassword(password)
	return u
}

func AllAdmins() []*Admin {
	admins := make([]*Admin, 0)
	x.Find(&admins)
	return admins
}

func InitAdmin() {
	total, _ := x.Count(new(Admin))
	if total <= 0 {
		pass := RandStringRunes(16)
		u := NewAdmin("admin", pass)
		x.Insert(u)
		log.Printf("no admin found, generate new admin [username: admin password: %s] \n", pass)
		ioutil.WriteFile("default_pass", []byte(fmt.Sprintf("username: admin\npassword; %s\n", pass)), 0600)
	}
}

func (u *Admin) GetnerateSalt() string {
	content := fmt.Sprintf("%s%d", u.Username, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", hash)
}

func (u *Admin) SetPassword(passwd string) error {
	salt := u.GetnerateSalt()
	saltedPassword := []byte(passwd + salt)
	enc, err := bcrypt.GenerateFromPassword(saltedPassword, bcryptCost)
	if err != nil {
		return err
	}
	u.EncryptedPassword = string(enc)
	u.Salt = salt
	return nil
}

func (u *Admin) CheckPassword(passwd string) error {
	saltedPassword := []byte(passwd + u.Salt)
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), saltedPassword)
}

func (u *Admin) BeforeInsert() {
	if u.Password != "" {
		u.SetPassword(u.Password)
		u.Password = ""
	}
}

func (u *Admin) BeforeUpdate() {
	if u.Password != "" {
		u.SetPassword(u.Password)
		u.Password = ""
	}

}
