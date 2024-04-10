package orm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/db_connection"
	"log"
	"time"
	"xorm.io/xorm"
)

type service struct {
	config db_connection.Config
}

var _ Service = (*service)(nil)

func NewService(c db_connection.Config) *service {
	return &service{config: c}
}

func (s *service) Init() *xorm.Engine {
	engine, err := xorm.NewEngine(s.config.DbDriver, fmt.Sprintf(s.config.DbDns, s.config.DbUser,
		s.config.DbPassword, s.config.DbHost, s.config.DbPort, s.config.DbName))
	if err != nil {
		log.Fatal(err)
	}
	engine.SetMaxIdleConns(s.config.MaxIdleCons)
	engine.SetMaxOpenConns(s.config.MaxOpenCons)
	engine.SetConnMaxLifetime(time.Minute * time.Duration(s.config.ConnMaxLifetime))
	return engine
}
