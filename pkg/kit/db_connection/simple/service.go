package simple

import (
	"database/sql"
	"fmt"
	_ "github.com/godror/godror"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/db_connection"
	"time"
)

type service struct {
	config db_connection.Config
}

var _ Service = (*service)(nil)

func NewService(cfg db_connection.Config) *service {
	return &service{
		config: cfg,
	}
}

func (s service) Init() *sql.DB {
	connLine := fmt.Sprintf(s.config.DbDns, s.config.DbUser, s.config.DbPassword,
		s.config.DbHost, s.config.DbPort, s.config.DbName)
	db, err := sql.Open("godror", connLine)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}
	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("error pinging db: %w", err))
	}
	db.SetMaxOpenConns(s.config.MaxOpenCons)
	db.SetMaxIdleConns(s.config.MaxIdleCons)
	db.SetConnMaxLifetime(time.Second * time.Duration(s.config.ConnMaxLifetime))
	return db
}
