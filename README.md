# WL Go Logger
It can help you to handle the file and log rotate and remove log rotated file issues.

-   [Installation](#installation)
-   [Quick Example](#quick-example)


## Installation
```bash
$ go get -u github.com/weilun-shrimp/wlgologger
```

## Quick Example
```golang
package main

import (
	"github.com/weilun-shrimp/wlgologger"
    "log"
    "sync"
)

func main() {
    logger := wlgologger.Logger{
        // A dir path that is where your log file location.
        // eg. /var/log/my_app
        DirPath: "/var/log/my_app",
        DirMode: os.FileMode(0644), // log dir permission. eg. 0644

        // A dir path that you want to put your rotated log files.
        // eg. /var/log/my_app/rotated
        RotateDirPath:  "/var/log/my_app/rotated",
        RotateDirMode:  os.FileMode(0644), // rotate's log dir permission. eg. 0644

        // It will add Logger's FileExtension automatically.
        // eg. "app_log" is your filename. And "log" is your file extension.
        // - Logger will create or search a file that file full name is "app_log.log"
        FileName:      "app",
        FileExtension: "log",                               // eg. "log" or "txt"
        FileFlag:      os.O_CREATE|os.O_WRONLY|os.O_APPEND, // (https://pkg.go.dev/os#pkg-constants)
        FileMode:      os.FileMode(0644),                   // log file permission. eg. 0644

        Prefix: "",            // log prefix
        Flag:   log.LstdFlags, // log flag (https://pkg.go.dev/log#pkg-constants)

        lock: sync.RWMutex{}
    }

    if err := logger.Init(); err != nil {
        fmt.Println("logger init fail")
        fmt.Println(err)
        return
    }
    defer func() {
        if err := logger.Release(); err != nil {
            fmt.Println("release logger fail")
            fmt.Println(err)
            return
        }
    }()

    // start to log content
    // Same as log.Println
    logger.Println("What ever you want to log.")

    // rotate your log file
    // This will move the log file to rotate dir and give it time.RFC3339 suffix. And create a new log file in the log dir.
    // eg. "app.log" is your log file. Func will move it into rotate dir and rename it like "app-2025-01-01T00:00:00Z07:00.log"
    if err := logger.Rotate(); err != nil {
        fmt.Println("Fail to rotate the log")
        fmt.Println(err)
        return
    }

    // log your content again after rotated. It will put content in the new log file.
    logger.Println("What ever you want to log.")

    // remove your rotated file
    // Func will remove all files that before the before_at argument
    // eg. If now is 2025-02-01 and there got two rotated files in your rotated dir called
    // - "app-2025-01-01T00:00:00Z07:00.log" and "app-2025-03-01T00:00:00Z07:00.log"
    // - Func will remove "app-2025-01-01T00:00:00Z07:00.log".
    if err := logger.RemoveRotated(time.Now()); err != nil {
        fmt.Println("Fail to remove rotated log")
        fmt.Println(err)
        return 
    }
}
```