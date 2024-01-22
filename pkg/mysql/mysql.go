package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// MysqlConfig is list config to create Mysql client
type MysqlConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// MysqlMethod is list all available method for Mysql
type MysqlMethod interface {
	GetDB() *gorm.DB
}

// Client is a wrapper for Mysql client
type Client struct {
	db *gorm.DB
}

// NewMysqlClient is func to create Mysql client
func NewMysqlClient(config interface{}) (MysqlMethod, error) {
	db, err := gorm.Open("Mysql", config)
	if err != nil {
		return nil, err
	}

	return &Client{db: db}, nil
}

// GetDB is func to return database client
func (c *Client) GetDB() *gorm.DB {
	return c.db
}
