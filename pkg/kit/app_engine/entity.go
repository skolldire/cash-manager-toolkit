package app_engine

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/gorm"
	"net/http"
	"sync"
)

type Engine struct {
	App                 *App
	HttpClient          map[string]http.Client
	DBOrmConnections    map[string]*gorm.DB
	DBSimpleConnections map[string]*sql.DB
	RepositoriesConfig  map[string]interface{}
	UsesCasesConfig     map[string]interface{}
	HandlerConfig       map[string]interface{}
}

type Service interface {
	Init() (Engine, error)
	Run() error
}

type App struct {
	Router *chi.Mux
	Port   int
	Scope  string
	mutex  sync.Mutex
}
