package app_engine

var _ Service = (*Engine)(nil)

func NewService() *Engine {
	return &Engine{
		App:                 nil,
		HttpClient:          nil,
		DBOrmConnections:    nil,
		DBSimpleConnections: nil,
		RepositoriesConfig:  nil,
		UsesCasesConfig:     nil,
		HandlerConfig:       nil,
	}
}

func (e Engine) Init() (Engine, error) {
	panic("implement me")
}

func (e Engine) Run() {
	panic("implement me")
}
