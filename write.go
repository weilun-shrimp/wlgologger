package wlgologger

import (
	"log"
)

// Same as log.Println
func (l *Logger) Println(v ...any) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	log.Println(v...)
}
