package gorm

import (
	stdLog "log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Options struct {
	Debug       bool
	Source      string
	OpenConn    int
	Idle        int
	IdleTimeout time.Duration
}

func New(gormConf *gorm.Config, o Options) *gorm.DB {
	if o.Debug {
		// 终端打印输入 sql 执行记录
		gormConf.Logger = gormLogger.New(
			stdLog.New(os.Stdout, "\r\n", stdLog.LstdFlags), // io writer
			gormLogger.Config{
				SlowThreshold:             2 * time.Second, // 慢查询 SQL 阈值
				Colorful:                  true,            // 禁用彩色打印
				IgnoreRecordNotFoundError: true,            // 忽略未找到记录错误
				LogLevel:                  gormLogger.Info, // Log lever
			},
		)
	} else {
		// 终端打印输入 sql 执行记录
		gormConf.Logger = gormLogger.New(
			stdLog.New(os.Stdout, "\r\n", stdLog.LstdFlags), // io writer
			gormLogger.Config{
				SlowThreshold:             2 * time.Second, // 慢查询 SQL 阈值
				Colorful:                  true,            // 禁用彩色打印
				IgnoreRecordNotFoundError: true,            // 忽略未找到记录错误
				LogLevel:                  gormLogger.Warn, // Log lever
			},
		)
	}
	stdLog.Println("opening connection to mysql")
	db, err := gorm.Open(mysql.Open(o.Source), gormConf)
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		stdLog.Fatalf("failed get connection for db: %v", err)
	}
	sqlDB.SetMaxIdleConns(o.Idle) // 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetConnMaxIdleTime(o.IdleTimeout)
	sqlDB.SetMaxOpenConns(o.OpenConn) // 设置打开数据库连接的最大数量。
	return db
}

func NewWithConf(o Options) *gorm.DB {
	if o.OpenConn == 0 {
		o.OpenConn = 50
	}
	if o.Idle == 0 {
		o.Idle = 10
	}
	if o.IdleTimeout == 0 {
		o.IdleTimeout = 14400 * time.Second
	}
	gormConf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名是否加 s
			TablePrefix:   "",
		},
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		DisableAutomaticPing: false,
		PrepareStmt:          false,
	}
	return New(gormConf, o)
}
