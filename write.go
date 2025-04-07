package wlgologger

// Same as log.Println
func (l *Logger) Println(v ...any) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.logger.Println(v...)
}
