package util

import (
	"os"
)

func CreateDir(a string) (err error) {
	_, err = os.Stat(a)
	if os.IsNotExist(err) {
		os.MkdirAll(a, os.ModePerm)
	}
	return err
}
