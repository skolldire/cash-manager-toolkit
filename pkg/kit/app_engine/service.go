package app_engine

var _ Service = (*Engine)(nil)

func NewService() *Engine {
	return &Engine{}
}

func (e Engine) Init() (Engine, error) {
	panic("implement me")
}

func (e Engine) Run() error {
	panic("implement me")
}
