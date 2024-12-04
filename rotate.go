package wlgologger

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"
)

// This will move the log file to rotate dir and give it time.RFC3339 suffix. And create a new log file in the log dir.
//
// eg. "app.log" is your log file. Func will move it into rotate dir and rename it like "app-2025-01-01T00:00:00Z07:00.log"
func (l *Logger) Rotate() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.file == nil {
		return ErrLoggerShouldBeInitializedFirst
	}

	// check log rotate store path
	if err := checkDir(l.RotateDirPath, l.RotateDirMode); err != nil {
		return err
	}

	if err := exec.Command(
		"mv",
		fmt.Sprintf("%s/%s.%s", l.DirPath, l.FileName, l.FileExtension),
		fmt.Sprintf("%s/%s-%s.%s", l.RotateDirPath, l.FileName, time.Now().Format(time.RFC3339), l.FileExtension),
	).Run(); err != nil {
		return err
	}

	// reset the logger
	if err := l.release(); err != nil {
		return err
	}
	if err := l.init(); err != nil {
		return err
	}

	return nil
}

// Func will remove all files that before the before_at argument
//
// eg. If now is 2025-02-01 and there got two rotated files in your rotated dir called
//
// - "app-2025-01-01T00:00:00Z07:00.log" and "app-2025-03-01T00:00:00Z07:00.log"
//
// - Func will remove "app-2025-01-01T00:00:00Z07:00.log".
func (l *Logger) RemoveRotated(before_at time.Time) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	// check log rotate store path
	if err := checkDir(l.RotateDirPath, l.RotateDirMode); err != nil {
		return err
	}

	dir_entry, err := os.ReadDir(l.RotateDirPath)
	if err != nil {
		return err
	}
	for _, file := range dir_entry {
		created_at_decoded := regexp.MustCompile(
			fmt.Sprintf("%s-(.*).%s$", l.FileName, l.FileExtension),
		).FindStringSubmatch(file.Name())
		if len(created_at_decoded) < 2 { // without create at (time.RFC3339 suffix)
			continue
		}
		created_at, err := time.Parse(time.RFC3339, created_at_decoded[1])
		if err != nil { // create at suffix format not time.RFC3339
			continue
		}
		if !created_at.Before(before_at) {
			continue
		}
		if err = os.Remove(fmt.Sprintf("%s/%s", l.RotateDirPath, file.Name())); err != nil {
			return err
		}
	}

	return nil
}
