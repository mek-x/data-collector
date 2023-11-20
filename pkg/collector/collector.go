package collector

import (
	"gitlab.com/mek_x/data-collector/pkg/parser"
)

type Collector interface {
	Start() error
	AddDataSource(path string, parser parser.Parser) error
	End()
}

type Init func(params any) Collector
type registeredCollectors map[string]Init

/* Main registry of all available collector classes */
var Registry = make(registeredCollectors)

func (r registeredCollectors) Add(name string, constructor Init) {
	r[name] = constructor
}
