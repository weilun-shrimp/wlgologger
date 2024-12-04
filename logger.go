package wlgologger

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"sync"
)

type Logger struct {
	// A dir path that is where your log file location.
	//
	// eg. /var/log/my_app
	DirPath string
	DirMode fs.FileMode // log dir permission. eg. 0644

	// A dir path that you want to put your rotated log files.
	//
	// eg. /var/log/my_app/rotated
	RotateDirPath string
	RotateDirMode fs.FileMode // rotate's log dir permission. eg. 0644

	// It will add Logger's FileExtension automatically.
	//
	// eg. "app_log" is your filename. And "log" is your file extension.
	//
	// - Logger will create or search a file that file full name is "app_log.log"
	FileName      string
	FileExtension string      // eg. "log" or "txt"
	FileFlag      int         // eg. os.O_CREATE|os.O_WRONLY|os.O_APPEND (https://pkg.go.dev/os#pkg-constants)
	FileMode      fs.FileMode // log file permission. eg. 0644

	Prefix string // log prefix
	Flag   int    // log flag (https://pkg.go.dev/log#pkg-constants)

	lock   sync.RWMutex
	file   *os.File
	logger *log.Logger
}

func (l *Logger) Init() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	// check log rotate dir path
	if err := checkDir(l.RotateDirPath, l.RotateDirMode); err != nil {
		return err
	}

	if err := l.init(); err != nil {
		return err
	}

	return nil
}

func (l *Logger) init() error {
	// check log dir path
	if err := checkDir(l.DirPath, l.DirMode); err != nil {
		return err
	}

	log_file, err := os.OpenFile(fmt.Sprintf("%s/%s.%s", l.DirPath, l.FileName, l.FileExtension), l.FileFlag, l.FileMode)
	if err != nil {
		return err
	}
	l.file = log_file
	if l.Flag == 0 {
		l.Flag = log.LstdFlags
	}
	l.logger = log.New(l.file, l.Prefix, l.Flag)
	return nil
}

func (l *Logger) Release() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if err := l.release(); err != nil {
		return err
	}

	return nil
}

func (l *Logger) release() error {
	if l.file != nil {
		if err := l.file.Close(); err != nil {
			return err
		}
	}
	if l.logger != nil {
		l.logger = nil
	}

	return nil
}
