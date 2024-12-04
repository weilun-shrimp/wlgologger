package wlgologger

import (
	"io/fs"
	"os"
)

// check the store dir path is valid.
//
// It will make a new dir if the dir is not exists.
func checkDir(path string, mode fs.FileMode) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, mode)
		if err != nil {
			return err
		}
	}

	return nil
}
