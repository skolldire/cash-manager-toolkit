package db_connection

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

type service struct {
	config Config
}

var _ Service = (*service)(nil)

func NewService(c Config) *service {
	return &service{config: c}
}

func (s *service) Init() *gorm.DB {
	db, err := gorm.Open("mysql", fmt.Sprintf(s.config.DbDns, s.config.DbUser,
		s.config.DbPassword, s.config.DbHost, s.config.DbPort, s.config.DbName))
	if err != nil {
		log.Fatal(err)
	}
	return db
}
