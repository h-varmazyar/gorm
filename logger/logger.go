package logger

import (
	"gorm.io/gorm/logger"
)

type Logger logger.Interface
type Writer logger.Writer
type Loglevel int

const (
	Silent Loglevel = iota + 1
	Error
	Warn
	Info
)