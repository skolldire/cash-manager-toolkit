package load_properties

type LoadProperties interface {
	Apply() (Config, error)
}

type AppConfig struct {
	Port     int    `json:"port"`
	LogLevel string `json:"log_level"`
	Name     string `json:"name"`
}

type Config struct {
	Application        AppConfig              `json:"application_config"`
	RepositoriesConfig map[string]interface{} `json:"repositories_config"`
	UsesCasesConfig    map[string]interface{} `json:"uses_cases_config"`
	HandlerConfig      map[string]interface{} `json:"handlers_config"`
}
