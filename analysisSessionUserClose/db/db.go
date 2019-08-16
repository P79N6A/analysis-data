package db

import (
	"fmt"

	"github.com/analysis-data/analysisSessionUserClose/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type GormInterface struct {
	gormDB *gorm.DB
}

func RegisterDB() (*GormInterface, error) {
	// (可选)设置最大空闲连接
	maxIdle := conf.AppConfig.DefaultInt("maxIdle", 60)
	// (可选) 设置最大数据库连接 (go >= 1.2)
	maxConn := conf.AppConfig.DefaultInt("maxConn", 300)
	dbLink := conf.AppConfig.DefaultString("link", "")
	if dbLink == "" {
		return nil, fmt.Errorf("no db link")
	}

	db, err := gorm.Open("mysql", dbLink)
	if err != nil {
		return nil, fmt.Errorf("RegisterDateBase error: " + err.Error())
	}

	db.DB().SetMaxIdleConns(maxIdle)
	db.DB().SetMaxOpenConns(maxConn)

	return &GormInterface{gormDB: db}, nil
}

func (db *GormInterface) RegisterTable(modules ...interface{}) error {
	db.gormDB.SingularTable(true)
	for _, module := range modules {
		if db.gormDB.HasTable(module) {
			continue
		}

		db.gormDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(module)
	}
	return nil
}

func (db *GormInterface) QuerySessions(startTime, endTime string, data interface{}) error {
	dbTemp := db.gormDB.Where("createTime >= ? and createTime <= ?", startTime, endTime).Find(data)
	if dbTemp.Error != nil {
		return dbTemp.Error
	}
	return nil
}
