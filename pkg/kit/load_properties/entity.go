package load_properties

import (
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/app"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/db_connection"
)

type LoadProperties interface {
	Apply() (Config, error)
}

type Config struct {
	Application               app.Config                        `json:"application_config"`
	DBOrmConnectionsConfig    []map[string]db_connection.Config `json:"db_orm_connections"`
	DBSimpleConnectionsConfig []map[string]db_connection.Config `json:"db_simple_connections"`
	RepositoriesConfig        map[string]interface{}            `json:"repositories_config"`
	UsesCasesConfig           map[string]interface{}            `json:"uses_cases_config"`
	HandlerConfig             map[string]interface{}            `json:"handlers_config"`
}

// []map[string]orm.Config
