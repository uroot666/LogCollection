package Log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// 创建全局的zap对象，后面通过一个函数提供给其它模块使用
var sugarLogger *zap.SugaredLogger

// 初始化zap对象，加载定义的一些配置
func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

// 修改时间编码器，在日志文件中使用大写字母记录日志级别
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 在zap中加入Lumberjack支持
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./Log/test.log", // 日志文件的位置
		MaxSize:    1,                // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 5,                // 保留旧文件的最大个数
		MaxAge:     30,               // 保留旧文件的最大天数
		Compress:   false,            // 是否压缩/归档旧文件
	}
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	// return zapcore.AddSync(lumberJackLogger)
}

// 对外提供一个接口，是外部模块可以调用创建好的log对象来记录日志
func GetLogObj() (*zap.SugaredLogger, error) {
	return sugarLogger, nil
}
