package wlgologger

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)


var log_file *os.File

const base_filename = "app"

func Init() {
	if os.Getenv("APP_LOG_PATH") == "" {
		panic(".env APP_LOG_PATH parameter is empty.")
	}

	// check log store path
	if _, err := os.Stat(os.Getenv("APP_LOG_PATH")); os.IsNotExist(err) {
		err := os.MkdirAll(os.Getenv("APP_LOG_PATH"), 0644)
		if err != nil {
			panic("Unable to create log store folder. Error msg: " + err.Error())
		}
	}

	// check log rotate store path
	if _, err := os.Stat(os.Getenv("APP_LOG_ROTATE_PATH")); os.IsNotExist(err) {
		err := os.MkdirAll(os.Getenv("APP_LOG_ROTATE_PATH"), 0644)
		if err != nil {
			panic("Unable to create log rotate store folder. Error msg: " + err.Error())
		}
	}

	var err error
	log_file, err = os.OpenFile(os.Getenv("APP_LOG_PATH")+"/"+base_filename+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("Open log file fail. Log file path: " + os.Getenv("APP_LOG_PATH") + " | error msg: " + err.Error())
	}

	log.SetOutput(log_file)
}

func Rotate() {
	if log_file == nil {
		log.Println("Rotate should be called after log init. Log rotate failed.")
		return
	}

	// check log rotate store path
	if _, err := os.Stat(os.Getenv("APP_LOG_ROTATE_PATH")); os.IsNotExist(err) {
		err := os.MkdirAll(os.Getenv("APP_LOG_ROTATE_PATH"), 0644)
		if err != nil {
			log.Println("Unable to create log rotate store folder. Error msg: " + err.Error())
			return
		}
	}

	exec.Command(
		"mv",
		os.Getenv("APP_LOG_PATH")+"/"+base_filename+".log",
		os.Getenv("APP_LOG_ROTATE_PATH")+"/"+base_filename+"-"+time.Now().Format(time.RFC3339)+".log",
	).Run()
	pre_log_file := log_file
	Init() // will reset log_file
	pre_log_file.Close()
}

func RotateRemove() error {
	rotate_trash_days, err := strconv.Atoi(os.Getenv("APP_LOG_ROTATE_TRASH_DAYS"))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	check_at := time.Now().AddDate(0, 0, 0-rotate_trash_days)
	dir_entry, err := os.ReadDir(os.Getenv("APP_LOG_ROTATE_PATH"))
	if err != nil {
		log.Println("Unable to read log rotate store folder. Error msg: " + err.Error())
		return err
	}
	for _, file := range dir_entry {
		created_at_decoded := regexp.MustCompile("app" + "-(.*).log$").FindStringSubmatch(file.Name())
		if len(created_at_decoded) < 2 {
			fmt.Println(file.Name() + " no created at.")
			continue
		}
		created_at, err := time.Parse(time.RFC3339, created_at_decoded[1])
		if err != nil {
			fmt.Println(file.Name() + " created at format is not time.RFC3339.")
			continue
		}
		if !created_at.Before(check_at) {
			continue
		}
		if err = os.Remove(os.Getenv("APP_LOG_ROTATE_PATH") + "/" + file.Name()); err != nil {
			Write("error", map[string]string{
				"from": "log_rotate_remove",
				"type": "fail_to_remove_file",
				"path": os.Getenv("APP_LOG_ROTATE_PATH") + "/" + file.Name(),
				"msg":  err.Error(),
			})
			return err
		}
	}
	return nil
}

func Close() {
	log_file.Close()
}