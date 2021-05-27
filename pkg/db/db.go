package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ErrorNotFound = errors.New("Not found")

type Connection struct {
	*gorm.DB
}

type ConnectionDetails struct {
	Host     string
	Port     int
	User     string
	Database string
	Password string
}

func NewConnection(details ConnectionDetails) (*Connection, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Warn, // Log level
			Colorful:      false,       // Disable color
		},
	)

	var err error
	c, err := gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		details.Host, details.Port, details.User, details.Database, details.Password)), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	return &Connection{c}, err
}

func (c *Connection) DoMigrate() error {
	err := c.AutoMigrate(
		&User{},
		&Display{},
		&NewsItem{},
		&NewsLanguageItem{},
		&TimeFrame{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) GetID(obj interface{}, id uint) error {
	res := c.First(obj, id)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrorNotFound
	}

	return nil
}

func (c *Connection) GetWhereIs(obj interface{}, property string, where interface{}) error {
	res := c.First(obj, fmt.Sprintf("%s = ?", property), where)

	if res.RowsAffected == 0 || res.Error == gorm.ErrRecordNotFound {
		return ErrorNotFound
	}

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (c *Connection) GetAllWhereIs(obj interface{}, property string, where interface{}) error {
	res := c.Find(obj, fmt.Sprintf("%s = ?", property), where)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrorNotFound
	}

	return nil
}
