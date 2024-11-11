package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


func Must(service string) *zap.SugaredLogger {

	config := zap.NewProductionConfig()

	config.EncoderConfig.EncodeTime = formattedTimeEncoder  // 时间戳格式
	config.DisableStacktrace = true                         // 打印堆栈
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)  // 日志级别, 过滤
	config.Encoding = "console"                             // 输出编码, 终端 console, 文件 json
	config.OutputPaths = []string{"stdout", "gua.log"}      // 输出路径
	config.InitialFields = map[string]any{
		"service": service,
	}

	baseLogger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return baseLogger.Sugar()
}

// 格式化时间戳格式   
func formattedTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
