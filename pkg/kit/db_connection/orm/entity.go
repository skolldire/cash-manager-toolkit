package orm

import (
	"xorm.io/xorm"
)

type Service interface {
	Init() *xorm.Engine
}

type Config struct {
	DbName     string `json:"db_name"`
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
	DbHost     string `json:"db_host"`
	DbPort     string `json:"db_port"`
	DbDns      string `json:"db_dns"`
}
