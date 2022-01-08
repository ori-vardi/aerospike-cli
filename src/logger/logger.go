package logger

import (
	"fmt"
	prop "github.com/magiconair/properties"
	"io"
	"log"
	"os"
)

type logLevel int

const (
	logLevelDebug logLevel = iota
	logLevelInfo
	logLevelWarn
	logLevelError
	logLevelFatal
)

const (
	PropertyLoggerCleanOldLog = "logger.clean.old.log"
	PropertyLoggerLevelFile   = "logger.level.file"
	PropertyLoggerPath        = "logger.path"
)

var (
	Debug   *Logger
	Info    *Logger
	Warning *Logger
	Error   *Logger
	Fatal   *Logger
)

func Init(props *prop.Properties) {
	logPath := props.MustGet(PropertyLoggerPath)

	if props.MustGetBool(PropertyLoggerCleanOldLog) {
		log.Printf("truncate (old) log file!, path: %s", logPath)
		ResetLogFile(logPath)
	}
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("error while init logger: ", err)
	} else {
		fmt.Printf("Writing log to: %s \n\n", logFile.Name())
		initialize(logFile, props)
	}
}

func convertToLogLevel(str string) logLevel {
	switch str {
	case "DEBUG":
		return logLevelDebug
	case "INFO":
		return logLevelInfo
	case "WARN":
		return logLevelWarn
	case "FATAL":
		return logLevelFatal
	case "ERROR":
		return logLevelError
	default:
		return logLevelInfo
	}
}

func getLogger(curLogLevel logLevel, fileLogLevel logLevel, logFile *os.File) io.Writer {

	var writers []io.Writer

	if curLogLevel >= fileLogLevel && logFile != nil {
		writers = append(writers, logFile)
	}

	return io.MultiWriter(writers...)
}

func initialize(logFile *os.File, props *prop.Properties) {
	fileLevel := convertToLogLevel(props.MustGet(PropertyLoggerLevelFile))

	infoLogger := getLogger(logLevelInfo, fileLevel, logFile)
	Info = new(infoLogger, "INFO: ", log.Ldate|log.Ltime)

	warnLogger := getLogger(logLevelWarn, fileLevel, logFile)
	Warning = new(warnLogger, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)

	errLogger := getLogger(logLevelError, fileLevel, logFile)
	Error = new(errLogger, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	debugLogger := getLogger(logLevelDebug, fileLevel, logFile)
	Debug = new(debugLogger, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	fatalLogger := getLogger(logLevelFatal, fileLevel, logFile)
	Fatal = new(fatalLogger, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func ResetLogFile(agentLogPath string) error {
	return os.Truncate(agentLogPath, 0)
}

type Logger struct {
	*log.Logger
}

func new(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{
		Logger: log.New(out, prefix, flag),
	}
}

func (log *Logger) Printf(format string, v ...interface{}) {
	log.Logger.Output(2, fmt.Sprintf(format, v...))
}

func (log *Logger) Println(v ...interface{}) {
	log.Logger.Output(2, fmt.Sprint(v...))
}
