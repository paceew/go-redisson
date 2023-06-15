package orm

import (
	"context"
	"time"

	"github.com/paceew/go-redisson/pkg/config"

	"github.com/aiwuTech/fileLogger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type GormConfig struct {
	MysqlConfig    *mysql.Config
	GormConfig     *gorm.Config
	ConnPoolConfig GormConnPoolConfig
}

type GormConnPoolConfig struct {
	MaxIdleConns int
	MaxOpenConns int
}

//	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
//		NamingStrategy: schema.NamingStrategy{
//		  TablePrefix: "t_",   // table name prefix, table for `User` would be `t_users`
//		  SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
//		  NoLowerCase: true, // skip the snake_casing of names
//		  NameReplacer: strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
//		},
//	  })
func AutoInitGormWithConfig() {
	level := fileLogger.LEVEL(config.VipCfg().GetIntWithDefault("database.logger.level", 1))
	if level > fileLogger.OFF {
		level = fileLogger.INFO
	}
	glv := 5 - int(level)
	if glv > 5 {
		glv = 4
	}

	glcfg := GormLogConfig{
		Level:                     glog.LogLevel(glv),
		IgnoreRecordNotFoundError: config.VipCfg().GetBoolWithDefault("database.logger.ignore_record_not_found", true),
		SlowThreshold:             config.VipCfg().GetDurationWithDefault("database.logger.slow_threshold", 200) * time.Millisecond,
		LogPath:                   config.VipCfg().GetStringWithDefault("database.logger.log_path", "./"),
		LogName:                   config.VipCfg().GetStringWithDefault("database.logger.log_name", "gorm.log"),
	}
	glogger := NewGormLogger(glcfg)

	dbcfg := config.VipCfg().GetStringMapString("database")
	for k := range dbcfg {
		if k == "logger" {
			continue
		}
		cfgprefix := "database." + k
		dsn := config.VipCfg().GetString(cfgprefix + ".dsn")
		if dsn == "" {
			continue
		}

		mysqlcfg := &mysql.Config{
			DSN:                       dsn,   // "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // DSN data source name
			SkipInitializeWithVersion: false, // 根据版本自动配置
		}
		gl := glogger.WithPrefix(k)
		gl.Info(context.TODO(), "%s GORM init", k)
		gcfg := &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   config.VipCfg().GetString(cfgprefix + ".table_prefix"),
				SingularTable: config.VipCfg().GetBoolWithDefault(cfgprefix+".singular_table", true),
			},
			Logger:                                   gl,
			DisableForeignKeyConstraintWhenMigrating: true,
			SkipDefaultTransaction:                   true,
		}
		gcpcfg := GormConnPoolConfig{
			MaxIdleConns: config.VipCfg().GetInt(cfgprefix + ".max_idle"),
			MaxOpenConns: config.VipCfg().GetInt(cfgprefix + ".max_open"),
		}
		cfg := GormConfig{
			MysqlConfig:    mysqlcfg,
			GormConfig:     gcfg,
			ConnPoolConfig: gcpcfg,
		}
		InitGorm(k, cfg)
	}
}
