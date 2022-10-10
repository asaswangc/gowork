package mysql_cli

import (
	"github.com/asaswangc/gowork/variable"
	"io"
	"log"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"gorm.io/gorm/logger"
)

func loggerFunc(cfg *LogCfg, runMode string) logger.Interface {
	// Sql日志文件配置
	var LogWriter = &lumberjack.Logger{
		MaxAge:     cfg.MaxAge,
		Filename:   cfg.LogPath,
		MaxSize:    cfg.Maxsize,
		Compress:   cfg.Compress,
		MaxBackups: cfg.MaxBackups,
	}

	// Sql日志配置
	var Logger = func() logger.Interface {
		var (
			IoWriter io.Writer
			LogLevel logger.LogLevel
		)
		switch runMode {
		case variable.ReleaseMode:
			IoWriter = LogWriter
			LogLevel = logger.Error
		default:
			IoWriter = os.Stdout
			LogLevel = logger.Info
		}
		return logger.New(log.New(IoWriter, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond, // 慢查询阀值
			LogLevel:                  LogLevel,               // 日志级别
			IgnoreRecordNotFoundError: false,                  // 是否忽略 RecordNotFoundError
			Colorful:                  false,                  // 是否启用日志颜色
		})
	}

	return Logger()
}
