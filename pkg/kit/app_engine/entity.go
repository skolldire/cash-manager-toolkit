package app_engine

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/app"
	"net/http"
)

type Engine struct {
	App                 *app.Service
	HttpClient          map[string]http.Client
	DBOrmConnections    map[string]*gorm.DB
	DBSimpleConnections map[string]*sql.DB
	RepositoriesConfig  map[string]interface{}
	UsesCasesConfig     map[string]interface{}
	HandlerConfig       map[string]interface{}
}

type Service interface {
	Init() (Engine, error)
	Run()
}
