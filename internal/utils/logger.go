package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type AppLogger struct {
	File   *os.File
	logger *log.Logger
}

func (appLogger *AppLogger) SetupLogger() {
	if DebugEnabled {
		appLogger.File = os.Stdout
	} else {
		getFilePath, err := setLogFile()
		if err != nil {
			log.Panic(err)
		}
		appLogger.File = getFilePath
	}
	appLogger.logger = log.New(appLogger.File, "Copy.io: ", log.LstdFlags)
}

func (appLogger *AppLogger) QuitLogger() {
	appLogger.File.Close()
}

func setLogFile() (*os.File, error) {

	var LOG_FILE string
	var path string

	path, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	// Crea il percorso della sottocartella "Copy IO" all'interno di "AppData/Local".
	// windows: C:\Users\<user>\AppData\Local\Copy IO\logs
	path = filepath.Join(path, "Copy IO", "logs")

	// read all logs files
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) > 0 {
		// sort them by last modified
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().Before(files[j].ModTime())
		})

		// check if last log file size exceeds 4MB
		if (files[len(files)-1].Size() / 1000) >= 4000 {
			// if yes, create a new log file
			t := time.Now()
			LOG_FILE = fmt.Sprintf("%s/%s.log", path, t.Format("02-01-2006 15h04m05s"))
		} else {
			LOG_FILE = fmt.Sprintf("%s/%s", path, files[len(files)-1].Name())
		}
	} else {
		t := time.Now()
		LOG_FILE = fmt.Sprintf("%s/%s.log", path, t.Format("02-01-2006 15h04m05s"))
	}

	// open log file
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return logFile, nil
}

func (appLogger *AppLogger) Debug(v ...any) {
	if DebugEnabled {
		appLogger.logger.Println("DEBUG: " + fmt.Sprint(v...))
	}
}

func (appLogger *AppLogger) Info(v ...any) {
	appLogger.logger.Println("INFO: " + fmt.Sprint(v...))
}

func (appLogger *AppLogger) Infof(format string, v ...any) {
	appLogger.logger.Printf("INFO: "+format, v...)
}

func (appLogger *AppLogger) Warning(v ...any) {
	appLogger.logger.Println("WARNING: " + fmt.Sprint(v...))
}

func (appLogger *AppLogger) Warningf(format string, v ...any) {
	appLogger.logger.Printf("WARNING: "+format, v...)
}

func (appLogger *AppLogger) Error(v ...any) {
	appLogger.logger.Println("ERROR: " + fmt.Sprint(v...))
}

func (appLogger *AppLogger) Errorf(format string, v ...any) {
	appLogger.logger.Printf("ERROR: "+format, v...)
}

func (appLogger *AppLogger) Fatal(v ...any) {
	appLogger.logger.Println("FATAL: " + fmt.Sprint(v...))
}

func (appLogger *AppLogger) Fatalf(format string, v ...any) {
	appLogger.logger.Printf("FATAL: "+format, v...)
}
