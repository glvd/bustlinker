package monitor

import "github.com/ipfs/go-ipfs/core/coreapi"

type monitor struct {
	api coreapi.CoreAPI
}

type Monitor monitor

type monitorData struct {
}

func New(repo string, cfg interface{}) (Monitor, error) {
	return Monitor{}, nil
}
