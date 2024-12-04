package wlgologger

import (
	"io/fs"
	"os"
)

// check the store dir path is valid.
//
// It will make a new dir if the dir is not exists.
func checkDir(path string, mode fs.FileMode) error {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(path, mode); err != nil {
			return err
		}
	}

	return err
}
