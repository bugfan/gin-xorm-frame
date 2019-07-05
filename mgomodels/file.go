package mgomodels

import (
	"io"
)

func InsertImage(filename string, data []byte) error {
	session := GetSession()
	defer session.Close()
	f, err := session.DB("sengine").GridFS("image").Create(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func GetImage(filename string, file io.Writer) (err error) {
	session := GetSession()
	defer session.Close()
	f, err := session.DB("sengine").GridFS("image").Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	io.Copy(file, f)
	return
}
