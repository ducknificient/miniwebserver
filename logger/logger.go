package logger

import (
	"encoding/json"
	"fmt"
	"miniwebserver/config"
	"os"
	"path/filepath"

	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(string)
	Error(string)
	Fatal(string)
	// SaveImage(file []byte)
}

type DefaultLogger struct {
	logger *zap.Logger
}

type LoggerData struct {
	Text  string
	Image []byte
	File  []byte
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func MyCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(filepath.Base(caller.FullPath()))
}

func NewLogger(config config.Configuration) (defaultlogger *DefaultLogger, err error) {

	defaultconfig := config.GetConfiguration()

	fmt.Println(*defaultconfig.PathLog)

	//check log folder is exist
	if _, err := os.Stat(*defaultconfig.PathLog); os.IsNotExist(err) {
		err := os.MkdirAll(*config.GetConfiguration().PathLog, 0766)
		if err != nil {
			return defaultlogger, err
		}
	}

	//create log file and setting rotate time (24 hours)
	// logFile := _pathlog + _filesep + "app-%Y-%m-%d-%H.log"
	logFile := *defaultconfig.PathLog + *defaultconfig.FileSep + "app-%Y-%m-%d.log"
	rotator, err := rotatelogs.New(
		logFile,
		rotatelogs.WithMaxAge(45*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour))
	if err != nil {
		return defaultlogger, err
	}

	// initialize the JSON encoding config
	encoderConfig := map[string]string{
		"levelEncoder": "lowercase",
		"levelKey":     "level",
		"timeKey":      "date",
		"timeEncoder":  "iso8601",
		"callerKey":    "caller",
		"messageKey":   "message",
	}
	data, _ := json.Marshal(encoderConfig)

	var encCfg zapcore.EncoderConfig
	if err := json.Unmarshal(data, &encCfg); err != nil {
		return defaultlogger, err
	}
	encCfg.EncodeTime = SyslogTimeEncoder
	encCfg.EncodeCaller = MyCaller

	// add the encoder config and rotator to create a new zap logger
	w := zapcore.AddSync(rotator)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encCfg),
		w,
		zap.InfoLevel)
	zap.New(core)

	defaultlogger = &DefaultLogger{
		logger: zap.New(core, zap.WithCaller(true), zap.AddStacktrace(zap.ErrorLevel)),
	}

	return defaultlogger, nil
}

func (l *DefaultLogger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *DefaultLogger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *DefaultLogger) Fatal(msg string) {
	l.logger.Fatal(msg)
}
