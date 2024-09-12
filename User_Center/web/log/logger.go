package log

import (
	"web/log/log_rotation"
	"web/models"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var SugarLogger *zap.SugaredLogger

func InitLogger() {

	encoder := getEncoder()
	writeSyncer := getLogWriter()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	// zap.AddCaller()  添加将调用函数信息记录到日志中的功能。
	logger := zap.New(core, zap.AddCaller())
	SugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	// 修改时间编码器
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 在日志文件中使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// NewConsoleEncoder 打印更符合人们观察的方式
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	log_rotation_logger := &log_rotation.Logger{
		Filename:   models.Overall_Situation_Logger.Logger.Filename,   // 文件位置
		MaxSize:    models.Overall_Situation_Logger.Logger.MaxSize,    // 进行切割之前,日志文件的最大大小(MB为单位)
		MaxSave:    models.Overall_Situation_Logger.Logger.MaxSave,    // 保留旧文件的最大天数
		MaxNumbers: models.Overall_Situation_Logger.Logger.MaxNumbers, // 保留旧文件的最大个数
		Compress:   models.Overall_Situation_Logger.Logger.Compress,   // 是否压缩/归档旧文件
	}
	return zapcore.AddSync(log_rotation_logger)
}
