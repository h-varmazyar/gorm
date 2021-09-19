package gorm

import (
	"fmt"
	myLogger "github.com/mrNobody95/gorm/logger"
	logrus "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

type DatabaseConfig struct {
	Type         DbType
	Username     string
	Password     string
	CharSet      string
	Name         string
	Host         string
	Port         int
	SSLMode      bool
	LogLevel  myLogger.Loglevel
	LogWriter    myLogger.Writer
	LogSlowThreshold time.Duration
	IgnoreRecordNotFoundError bool
	ColorfulLog bool
}

type DB struct {
	gorm.DB
}

type DbType string

const (
	MYSQL      DbType = "mysql"
	SQLITE     DbType = "sqlite"
	PostgreSQL DbType = "postgre"
)

func (conf *DatabaseConfig) Initialize(models ...interface{}) (error, *DB) {
	var db *gorm.DB
	var err error

	if conf.LogWriter == nil {
		conf.LogWriter = log.New(os.Stdout, "\r\n", log.LstdFlags)
	}
	if conf.LogSlowThreshold == 0 {
		conf.LogSlowThreshold=time.Second
	}
	newLogger := logger.New(
		conf.LogWriter,
		logger.Config{
			SlowThreshold:             conf.LogSlowThreshold,
			LogLevel:                  logger.LogLevel(conf.LogLevel),
			IgnoreRecordNotFoundError: conf.IgnoreRecordNotFoundError,
			Colorful:                  conf.ColorfulLog,
		},
	)

	switch conf.Type {
	case MYSQL:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Host, conf.Port, conf.Name)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			if strings.Contains(err.Error(), "Unknown database") {
				logrus.Info("creating database ", conf.Name)
				dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Host, conf.Port)
				db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
					Logger: newLogger,
				})
				if err == nil {
					db.Exec(fmt.Sprintf("CREATE DATABASE %s;", conf.Name))
					db.Exec(fmt.Sprintf("USE %s;", conf.Name))
				} else {
					return err, nil
				}
			} else {
				return err, nil
			}
		}
	case PostgreSQL:
		ssl := "disable"
		if conf.SSLMode {
			ssl = "enable"
		}
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC", conf.Host, conf.Username, conf.Password, conf.Name, conf.Port, ssl)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			return err, nil
		}
	case SQLITE:
		db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
			Logger: newLogger,
		})
	}

	return migrate(db, models), &DB{DB: *db}
}

func migrate(db *gorm.DB, models ...interface{}) error {
	for _, intArr := range models {
		for _, model := range intArr.([]interface{}) {
			if err := db.AutoMigrate(model); err != nil {
				return err
			}
		}
	}
	return nil
}
