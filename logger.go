package wlgologger

import (
	"io"
	"os"
)

type Logger struct {
	StorePath       string
	RotateStorePath string
	FileName        string
	file            *os.File
}

func test() {
	multiWriter := io.MultiWriter(file1, file2)
}
