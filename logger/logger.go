package logger

import "gorm.io/gorm/logger"

type Logger logger.Interface
type Config logger.Config
type Writer logger.Writer
type Loglevel logger.LogLevel


const (
	Silent Loglevel = iota + 1
	Error
	Warn
	Info
)