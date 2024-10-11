package app_engine

import (
	"database/sql"
	"github.com/skolldire/cash-manager-toolkit/pkg/client/log"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/app"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/db_connection"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/db_connection/orm"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/db_connection/simple"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/load_properties/viper"
	"github.com/skolldire/cash-manager-toolkit/pkg/server/tcp"
	"xorm.io/xorm"
)

func NewApp() *Engine {
	v := viper.NewService()
	c, err := v.Apply()
	if err != nil {
		panic(err)
	}
	return &Engine{
		App:                 app.NewService(c.Application),
		Tracer:              creteTracer(c.Application),
		HttpClient:          nil,
		DBOrmConnections:    createDBOrmConnections(c.DBOrmConnectionsConfig),
		DBSimpleConnections: createDBSimpleConnections(c.DBSimpleConnectionsConfig),
		RepositoriesConfig:  c.RepositoriesConfig,
		UsesCasesConfig:     c.UsesCasesConfig,
		HandlerConfig:       c.HandlerConfig,
	}
}

func creteTracer(config app.Config) log.Service {
	return log.NewService(config.LogLevel)
}

func createDBOrmConnections(cs []map[string]db_connection.Config) map[string]*xorm.Engine {
	connections := make(map[string]*xorm.Engine)
	for _, c := range cs {
		for k, v := range c {
			connections[k] = orm.NewService(v).Init()
		}
	}
	return connections
}

func createDBSimpleConnections(cs []map[string]db_connection.Config) map[string]*sql.DB {
	connections := make(map[string]*sql.DB)
	for _, c := range cs {
		for k, v := range c {
			connections[k] = simple.NewService(v).Init()
		}
	}
	return connections
}

func createTCPServer(cs []map[string]tcp.Config, log log.Service) map[string]*tcp.Service {
	servers := make(map[string]*tcp.Service)
	for _, c := range cs {
		for k, v := range c {
			d := tcp.Dependencies{
				Config: v,
				Log:    log,
			}
			servers[k] = tcp.NewService(d)
		}
	}
	return servers
}
