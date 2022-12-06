package log

import (
	"os"
	"strconv"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/rs/zerolog"
)

// logger - abstraction from the underlying logger impl
type logger struct {
	logger zerolog.Logger
}

func (l logger) Trace(msg string, args ...interface{}) {
	l.write(l.logger.Trace(), msg, args)
}

func (l logger) Debug(msg string, args ...interface{}) {
	l.write(l.logger.Debug(), msg, args)
}

func (l logger) Info(msg string, args ...interface{}) {
	l.write(l.logger.Info(), msg, args)
}

func (l logger) Warn(msg string, args ...interface{}) {
	l.write(l.logger.Warn(), msg, args)
}

func (l logger) Error(msg string, args ...interface{}) {
	l.write(l.logger.Error(), msg, args)
}

func (l *logger) Level() log.Level {
	return l.Level()
}

func (l logger) write(e *zerolog.Event, msg string, args []interface{}) {
	e.Fields(args)
	e.Msg(msg)
}

func Debug(msg string, args ...interface{}) {
	Logger().Debug(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	Logger().Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	Logger().Error(msg, args...)
}

// singleton
var inst *logger

func Logger() logger {
	if inst == nil {
		l := pluginLogger()
		if l != "loki" {
			zl := zerolog.New(os.Stderr).With().Timestamp().Logger()
			inst = &logger{
				logger: zl,
			}
			return *inst
		}

		lokiLogger := newLokiLogger()
		zl := zerolog.New(lokiLogger).Level(zerolog.Level(lokiLogger.Level))
		inst = &logger{
			logger: zl,
		}
		return *inst
	}
	return *inst
}

func pluginLogger() string {
	var envVar = "GF_PLUGIN_LOGGER"
	if v, ok := os.LookupEnv(envVar); ok {
		return v
	}
	return "console"
}

func newLokiLogger() *lokiLogger {
	return &lokiLogger{
		URL:        envVal("LOGGER_URL"),
		Key:        envVal("LOGGER_KEY"),
		BufferSize: envValInt("LOGGER_BUFFER"),
		Level:      envValInt("LOGGER_LEVEL"),
		Labels:     envVarLabels(),
	}
}

func envVal(key string) string {
	if v, ok := os.LookupEnv("GF_PLUGIN_" + key); ok {
		return v
	}
	// fmt.Println("missing env var: GF_PLUGIN_" + key)
	return ""
}

func envValInt(key string) int8 {
	if v, ok := os.LookupEnv("GF_PLUGIN_" + key); ok {
		if b, err := strconv.ParseInt(v, 10, 8); err == nil {
			return int8(b)
		}
		return 0
	}
	return 0
}

func envVarLabels() map[string]string {
	labels := map[string]string{}
	envLabels := envVal("LOGGER_LABELS")
	if envLabels != "" {
		labelList := strings.Split(envLabels, ",")
		for _, l := range labelList {
			labelValue := strings.Split(l, ":")
			if len(labelValue) > 1 {
				labels[labelValue[0]] = labelValue[1]
			}
		}
	}
	return labels
}

type Labels map[string]string
