package infrastructure

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"bracelet-ticket-system-be/pkg/xlogger"
)

func dbSetup() (*gorm.DB, error) {
	logger := xlogger.Logger
	l := gormLogger.Default.LogMode(gormLogger.Silent)

	if cfg.DbDriver == "mysql" {
		db, err := gorm.Open(mysql.New(mysql.Config{
			DSN: cfg.DbDsn,
		}), &gorm.Config{
			Logger: l,
		})
		if err != nil {
			logger.Info().Msg("Failed connect to database!")
			return nil, err
		}
		logger.Info().Msg("Successfully connected to the database!")
		return db, nil
	}
	return nil, fmt.Errorf("unsupported database driver")
}
