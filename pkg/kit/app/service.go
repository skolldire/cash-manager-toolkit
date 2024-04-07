package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/app/ping"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

var _ Service = (*App)(nil)

func NewService(c Config) *App {
	return &App{
		Router: initRoutes(),
		Port:   setPort(c.Port),
		Scope:  setScope(c.Scope),
	}
}

func (a App) Run() {
	err := http.ListenAndServe(a.Port, a.Router)
	if err != nil {
		panic(err)
	}
}

func initRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/ping", ping.NewService().Apply())
	r.Mount("/swagger", httpSwagger.WrapHandler)
	return r
}

func setPort(p string) string {
	if p == "" {
		return appDefaultPort
	}
	return p
}

func setScope(s string) string {
	if s == "" {
		return appDefaultScope
	}
	return s
}
