package app_engine

import (
	"database/sql"
	"github.com/skolldire/cash-manager-toolkit/pkg/client/log"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/app"
	"github.com/skolldire/cash-manager-toolkit/pkg/server/tcp"
	"net/http"
	"xorm.io/xorm"
)

type Engine struct {
	App                 *app.App
	Tracer              log.Service
	HttpClient          map[string]http.Client
	DBOrmConnections    map[string]*xorm.Engine
	DBSimpleConnections map[string]*sql.DB
	TCPServer           map[string]*tcp.Service
	RepositoriesConfig  map[string]interface{}
	UsesCasesConfig     map[string]interface{}
	HandlerConfig       map[string]interface{}
}
