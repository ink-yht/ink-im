package ioc

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"ink-im-server/internal/repository/dao"
	"ink-im-server/pkg/logger"
	"time"
)

func InitDB(l logger.Logger) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var c Config
	err := viper.UnmarshalKey("MySQL", &c)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败: %s \n", err))
	}
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			// 慢查询阈值，只有执行时间超过这个阈值，才会使用
			// 50ms,100ms
			SlowThreshold:             time.Millisecond * 50,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  glogger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}

//type gormLoggerFunc func(msg string, fields ...logger.Field)
//
//func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
//	g(msg, logger.Field{Key: "args",Value: args})
//}
