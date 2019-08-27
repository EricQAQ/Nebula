package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/EricQAQ/Nebula/config"
)

const (
	defaultLogTimeFormat = "2006-01-02 15:04:05.000"
	defaultLogFormat     = "text"
	defaultLogLevel      = log.InfoLevel
)

func stringToLogLevel(level string) log.Level {
	switch strings.ToLower(level) {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn", "warning":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	}
	return defaultLogLevel
}

func logTypeToColor(level log.Level) string {
	switch level {
	case log.DebugLevel:
		return "[0;37"
	case log.InfoLevel:
		return "[0;36"
	case log.WarnLevel:
		return "[0;33"
	case log.ErrorLevel:
		return "[0;31"
	case log.FatalLevel:
		return "[0;31"
	case log.PanicLevel:
		return "[0;31"
	}

	return "[0;37"
}

type nebulaLoggerFormatter struct {
	EnableColors  bool
	EnableSorting bool
}

// Format implements logrus.Formatter
// 2018-11-03 13:35:59.666 ~/this_is_a_fake_proj/xxx.go:67 [INFO] @@ this is the log format. field_1=eric field_2=zhang
// |------timestamp------| |----------file:line----------| [level]   |-------message-------| |---field, value pairs---|
func (lf *nebulaLoggerFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	fmt.Fprintf(b, "%s ", entry.Time.Format(defaultLogTimeFormat))
	if entry.Caller != nil {
		fmt.Fprintf(b, "%s:%d", entry.Caller.File, entry.Caller.Line)
	}
	if lf.EnableColors {
		colorStr := logTypeToColor(entry.Level)
		fmt.Fprintf(b, "\033%sm ", colorStr)
	}
	fmt.Fprintf(b, " [%s] @@ %s", entry.Level.String(), entry.Message)

	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	if lf.EnableSorting {
		sort.Strings(keys)
	}
	for _, k := range keys {
		fmt.Fprintf(b, " %v=%v", k, entry.Data[k])
	}

	b.WriteByte('\n')

	if lf.EnableColors {
		b.WriteString("\033[0m")
	}
	return b.Bytes(), nil
}

func createLoggerFormatter(format string) log.Formatter {
	switch strings.ToLower(format) {
	case "text":
		return &nebulaLoggerFormatter{}
	case "json":
		return &log.JSONFormatter{
			TimestampFormat: defaultLogTimeFormat,
		}
	case "console":
		return &log.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: defaultLogTimeFormat,
		}
	case "highlight":
		return &nebulaLoggerFormatter{
			EnableColors: true,
		}
	default:
		return &nebulaLoggerFormatter{}
	}
}

// initFileLog initializes file based logging options.
func initFileLog(cfg *config.LogConfig) error {
	var output io.Writer
	if !cfg.LogRotate {
		var err error
		output, err = os.OpenFile(cfg.File, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			if os.IsNotExist(err) {
				return errors.New("can't use directory as log file name")
			}
			return errors.New(err.Error())
		}
	} else {
		output = &lumberjack.Logger{
			Filename:   cfg.File,
			MaxSize:    int(cfg.MaxSize),
			MaxBackups: int(cfg.MaxBackups),
			MaxAge:     int(cfg.MaxDays),
			LocalTime:  true,
		}
	}
	log.SetOutput(output)
	return nil
}

func CreateLoggerFromConfig(cfg *config.NebulaConfig) error {
	logCfg := cfg.Log
	log.SetLevel(stringToLogLevel(logCfg.Level))

	if logCfg.Format == "" {
		logCfg.Format = defaultLogFormat
	}
	logFormatter := createLoggerFormatter(logCfg.Format)
	log.SetReportCaller(true)
	log.SetFormatter(logFormatter)

	if len(logCfg.File) != 0 {
		if err := initFileLog(logCfg); err != nil {
			return errors.Trace(err)
		}
	} else {
		log.SetOutput(os.Stdout)
	}
	return nil
}
