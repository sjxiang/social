package logger

import "testing"


// 测试日志输出
func TestLog(t *testing.T) {
	logger := Must("social")	
	
	logger.Info("测试")
	logger.Info("开发")
	logger.Debug("调试")
}

/*

logger := zap.Must(zap.NewProduction()).Sugar()
defer logger.Sync()

*/