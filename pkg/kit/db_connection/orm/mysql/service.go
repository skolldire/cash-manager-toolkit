package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/db_connection/orm"
	"log"
)

type service struct {
	config orm.Config
}

var _ orm.Service = (*service)(nil)

func NewService(c orm.Config) *service {
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
