package wlgologger

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestInit(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		panic("Load .env file error")
	}

	logger := &Logger{
		DirPath: os.Getenv("DIR_PATH"),
		DirMode: os.FileMode(0755),

		RotateDirPath: os.Getenv("ROTATED_DIR_PATH"),
		RotateDirMode: os.FileMode(0755),

		FileName:      os.Getenv("FILE_NAME"),
		FileExtension: os.Getenv("FILE_EXTENSION"),
		FileFlag:      os.O_CREATE | os.O_WRONLY | os.O_APPEND, // (https://pkg.go.dev/os#pkg-constants)
		FileMode:      os.FileMode(0755),

		Prefix: "",            // log prefix
		Flag:   log.LstdFlags, // log flag (https://pkg.go.dev/log#pkg-constants)
	}

	if err := logger.Init(); err != nil {
		fmt.Printf("%+v\n", logger)
		t.Fatalf("Init() error = %v", err)
	}
}
