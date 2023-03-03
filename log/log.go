package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/jdxj/oh-my-feed/config"
)

var (
	sugar *zap.SugaredLogger
)

func Init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	syncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:  config.Logger.Filename,
		MaxAge:    config.Logger.MaxAge,
		LocalTime: true,
	})
	levelEnabler := zap.LevelEnablerFunc(func(level zapcore.Level) bool { return level >= zapcore.Level(config.Logger.Level) })
	core := zapcore.NewCore(encoder, syncer, levelEnabler)

	sugar = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

func Debugf(template string, args ...any) {
	sugar.Debugf(template, args...)
}

func Infof(template string, args ...any) {
	sugar.Infof(template, args...)
}

func Warnf(template string, args ...any) {
	sugar.Warnf(template, args...)
}

func Errorf(template string, args ...any) {
	sugar.Errorf(template, args...)
}
