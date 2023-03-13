package postgres

import (
	"lab_2/config"

	extraClausePlugin "github.com/WinterYukky/gorm-extra-clause-plugin"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	Read  *gorm.DB
	Write *gorm.DB
}

func NewStorage(config config.SQLDatabase) (d *Storage, err error) {
	d = new(Storage)
	d.Read, err = d.initNewConnection(config.DSN)
	if err != nil {
		return
	}
	d.Write, err = d.initNewConnection(config.DSN)
	if err != nil {
		return
	}

	d.Write.AutoMigrate() //todo
	return
}

func (d *Storage) initNewConnection(dsn string) (*gorm.DB, error) {
	con, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(config.GetGormLogLevel()),
		// SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}
	con.Use(extraClausePlugin.New())
	return con, err
}
