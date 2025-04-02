package wlgologger

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestCheckDir(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		panic("Load .env file error")
	}

	err := checkDir(os.Getenv("ROTATED_DIR_PATH"), 0755)
	if err != nil {
		t.Errorf("checkDir() error = %v", err)
	}

	err = checkDir(os.Getenv("DIR_PATH"), 0755)
	if err != nil {
		t.Errorf("checkDir() error = %v", err)
	}

}
