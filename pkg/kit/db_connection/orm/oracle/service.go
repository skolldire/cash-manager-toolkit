package oracle

import (
	"fmt"
	_ "github.com/godror/godror"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/db_connection/orm"
	"log"
	"xorm.io/xorm"
)

type service struct {
	config orm.Config
}

var _ orm.Service = (*service)(nil)

func NewService(c orm.Config) *service {
	return &service{config: c}
}

func (s *service) Init() *xorm.Engine {
	engine, err := xorm.NewEngine("godror", fmt.Sprintf(s.config.DbDns, s.config.DbUser,
		s.config.DbPassword, s.config.DbHost, s.config.DbPort, s.config.DbName))
	if err != nil {
		log.Fatal(err)
	}
	return engine
}
