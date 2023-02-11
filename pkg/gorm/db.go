package gorm

import (
	"fmt"
	"time"

	"github.com/si-bas/go-rest-boilerplate/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func ConnectDB() *gorm.DB {

	// Write
	db := connectMysqlDb(config.Config.Db)

	// Read
	db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{mysql.Open(configToDsn(config.Config.Db))},
		Policy:   dbresolver.RandomPolicy{},
	}))

	return db
}

func connectMysqlDb(dbConfig config.DB) *gorm.DB {

	dsn := configToDsn(dbConfig)

	logLevel := logger.Silent
	if config.Config.App.Debug {
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		panic("error connecting to database, err=" + err.Error())
	}

	mysqlDb, _ := db.DB()
	mysqlDb.SetMaxOpenConns(dbConfig.Connection.Open)
	mysqlDb.SetMaxIdleConns(dbConfig.Connection.Idle)

	mili, _ := time.ParseDuration(fmt.Sprintf("%dms", dbConfig.Connection.TTL))
	mysqlDb.SetConnMaxLifetime(time.Duration(mili.Nanoseconds()))

	return db
}

func configToDsn(dbConfig config.DB) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
}
