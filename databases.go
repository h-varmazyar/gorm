package gorm

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

var err error
var db  *gorm.DB


type DatabaseConfig struct {
	Type     DbType
	Username string
	Password string
	CharSet  string
	Name     string
	Host     string
	Port     int
	SSLMode  bool
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
	switch conf.Type {
	case MYSQL:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Host, conf.Port, conf.Name)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			if strings.Contains(err.Error(), "Unknown database") {
				log.Info("creating database ", conf.Name)
				dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Host, conf.Port, "")
				db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
				if err == nil {
					err=db.Exec(fmt.Sprintf("create database %s;", conf.Name)).Error
				}
			}
			return err, nil
		}
	case PostgreSQL:
		ssl := "disable"
		if conf.SSLMode {
			ssl = "enable"
		}
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC", conf.Host, conf.Username, conf.Password, conf.Name, conf.Port, ssl)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return err, nil
		}
	case SQLITE:
		db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	}

	//return migrate(models), (*DB)(db)
	return migrate(models), &DB{DB:*db}
}

func migrate(models ...interface{}) error {
	for _, intArr := range models {
		//fmt.Println("type is:", )
		fmt.Println(reflect.TypeOf(intArr))
		for _, model := range intArr.([]interface{}) {
			if err := db.AutoMigrate(model); err != nil {
				return err
			}
		}

	}

	return nil
}
