package tcp

import (
	"github.com/skolldire/cash-manager-toolkit/pkg/client/log"
)

type Service interface {
	Init(f ProcessingFunc)
}

type ProcessingFunc func(msg string) (string, error)

type Dependencies struct {
	Port string
	Log  log.Service
}
