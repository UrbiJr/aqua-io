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

type Logger struct {
	File   *os.File
	logger *log.Logger
}

// And just go global.
var defaultLogger *Logger

func NewLogger() {
	defaultLogger = new(Logger)
	getFilePath, err := setLogFile()
	if err != nil {
		log.Panic(err)
	}
	defaultLogger.File = getFilePath
	defaultLogger.logger = log.New(getFilePath, "CACTUS-AIO: ", log.LstdFlags)
}

func QuitLogger() {
	defaultLogger.File.Close()
}

func setLogFile() (*os.File, error) {

	var LOG_FILE string
	var path string

	path, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	// Crea il percorso della sottocartella "NyxAIO" all'interno di "AppData/Local".
	// windows: C:\Users\<user>\AppData\Local\Roaming\NyxAIO\logs
	path = filepath.Join(path, "NyxAIO", "logs")

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

func Debug(v ...any) {
	if DebugEnabled {
		defaultLogger.logger.Println("DEBUG: " + fmt.Sprint(v...))
	}
}

func Info(v ...any) {
	defaultLogger.logger.Println("INFO: " + fmt.Sprint(v...))
}

func Infof(format string, v ...any) {
	defaultLogger.logger.Printf("INFO: "+format, v...)
}

func Warning(v ...any) {
	defaultLogger.logger.Println("WARNING: " + fmt.Sprint(v...))
}

func Warningf(format string, v ...any) {
	defaultLogger.logger.Printf("WARNING: "+format, v...)
}

func Error(v ...any) {
	defaultLogger.logger.Println("ERROR: " + fmt.Sprint(v...))
}

func Errorf(format string, v ...any) {
	defaultLogger.logger.Printf("ERROR: "+format, v...)
}

func Fatal(v ...any) {
	defaultLogger.logger.Println("FATAL: " + fmt.Sprint(v...))
}

func Fatalf(format string, v ...any) {
	defaultLogger.logger.Printf("FATAL: "+format, v...)
}
