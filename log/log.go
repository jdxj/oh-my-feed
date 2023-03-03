package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger *zap.SugaredLogger
)

func Init(writer *lumberjack.Logger) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	syncer := zapcore.AddSync(writer)
	levelEnabler := zap.LevelEnablerFunc(func(level zapcore.Level) bool { return level >= zapcore.DebugLevel })
	core := zapcore.NewCore(encoder, syncer, levelEnabler)
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0)).Sugar()
}
