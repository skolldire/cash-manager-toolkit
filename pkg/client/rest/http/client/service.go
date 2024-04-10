package client

import "net/http"

type service struct {
	configuration Configuration
}

var _ Service = (*service)(nil)

func NewService(c Configuration) *service {
	return &service{
		configuration: c,
	}
}

func (s service) Init() http.Client {
	panic("implement me")
}
