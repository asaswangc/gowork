package logger

import (
	"fmt"
	"github.com/asaswangc/gowork/variable"
	"log"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
)

// Init 创建自定义zap logger对象
func Init(cfg Cfg, hooks ...func(zapcore.Entry) error) {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	// 若为debug模式，创建debug日志级别的logger对象，直接输出到屏幕，不写入文件
	if variable.Global.Get(variable.RunMode) != variable.ReleaseMode {
		logger, err := zap.NewDevelopment(zap.Hooks(hooks...))
		if err != nil {
			fmt.Printf("创建zap日志包失败，详情：%s\n", err.Error())
		}
		Logger = logger
	}

	// 设置日志内容格式，以及日志输出格式。默认为人类可读格式；若配置了json，则输出为json格式
	encoderConf := genEncoderConf()
	encoder := zapcore.NewConsoleEncoder(encoderConf)
	if cfg.JsonEncoder {
		encoder = zapcore.NewJSONEncoder(encoderConf)
	}

	// 错误日志
	errWriter := &lumberjack.Logger{
		Filename:   cfg.ErrLog,
		MaxSize:    cfg.Maxsize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		Compress:   cfg.Compress,
	}

	// 警告日志
	warnWriter := &lumberjack.Logger{
		Filename:   cfg.WarnLog,
		MaxSize:    cfg.Maxsize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		Compress:   cfg.Compress,
	}
	// 普通日志
	infoWriter := &lumberjack.Logger{
		Filename:   cfg.InfoLog,
		MaxSize:    cfg.Maxsize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		Compress:   cfg.Compress,
	}

	// 日志级别配置，不能直接写zap.InfoLevel等，否则在写error级别的log时，info、warn也会写一份
	errLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zap.ErrorLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv > zap.InfoLevel && lv <= zap.WarnLevel
	})
	infoLevel := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv > zap.DebugLevel && lv <= zap.InfoLevel
	})

	// 启用多个输出流，不同级别的日志写到不同的日志文件中
	// 由于启用了多个输出流，所以配置文件中不必设置log_level，没有意义
	writers := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.AddSync(errWriter), errLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
	}

	// 创建zap logger对象，同时添加两个option：日志打印行号、自定义hook
	Logger = zap.New(zapcore.NewTee(writers...), zap.AddCaller(), zap.AddCallerSkip(1), zap.Hooks(hooks...))
}

// genEncoderConf 生成EncoderConfig，用于配置日志格式
func genEncoderConf() zapcore.EncoderConfig {
	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.EncodeTime = zapTimeEncoder               // 日志规范要求时间格式到毫秒
	encoderConf.TimeKey = "created_at"                    // 时间戳的key使用timestamp，根据model的定义进行设置
	encoderConf.MessageKey = "message"                    // 消息的key使用message
	encoderConf.EncodeLevel = zapcore.CapitalLevelEncoder // 日志规范要求日志级别为大写格式
	return encoderConf
}

// zapTimeEncoder 用于日志时间格式化，到毫秒级
func zapTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(variable.TimeFormat))
}
