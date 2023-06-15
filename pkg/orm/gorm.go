package orm

import (
	"context"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	gormMysqls map[string]*gorm.DB
)

var (
	ErrDBNotFound = errors.New("the db singletons not found or not initialization")
)

//GetGormSingleton 获取GetGormSingleton单例
func GetGormSingleton(name string) (gdb *gorm.DB, err error) {
	if name != "" && gormMysqls != nil {
		gdb = gormMysqls[name]
	}

	if gdb == nil {
		err = ErrDBNotFound
	}
	return
}

type GormModel interface {
	DbName() string
}

func GetGormSingletonByModel(gormmodel GormModel) (gdb *gorm.DB, err error) {
	return GetGormSingleton(gormmodel.DbName())
}

func setdb(name string, db *gorm.DB) {
	if gormMysqls == nil {
		gormMysqls = make(map[string]*gorm.DB)
	}
	gormMysqls[name] = db
}

func InitGorm(dbname string, cfg GormConfig) error {
	mysqlConfig := *cfg.MysqlConfig
	gconfig := cfg.GormConfig
	if db, err := gorm.Open(mysql.New(mysqlConfig), gconfig); err != nil {
		cfg.GormConfig.Logger.Error(context.TODO(), "default gorm open error:%s", err.Error())
		return err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(cfg.ConnPoolConfig.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.ConnPoolConfig.MaxOpenConns)
		setdb(dbname, db)
	}
	return nil
}
