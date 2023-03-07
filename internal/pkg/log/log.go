package log

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/jdxj/oh-my-feed/internal/pkg/config"
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
	if config.Logger.Filename == "" {
		syncer = zapcore.AddSync(os.Stdout)
	}

	levelEnabler := zap.LevelEnablerFunc(func(level zapcore.Level) bool { return level >= zapcore.Level(config.Logger.Level) })
	core := zapcore.NewCore(encoder, syncer, levelEnabler)

	sugar = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

func Desugar() *zap.Logger {
	return sugar.Desugar()
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

func Fatalf(template string, args ...any) {
	sugar.Fatalf(template, args...)
}

func Sync() {
	sugar.Sync()
}
