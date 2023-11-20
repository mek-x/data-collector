package parser

import "gitlab.com/mek_x/data-collector/pkg/datastore"

type Parser interface {
	Parse(buf []byte) error
	AddVar(name string, val string)
	AddConv(name string, val string)
}

type Init func(name string, ds datastore.DataStore) Parser
type registry map[string]Init

/* Main registry of all available collector classes */
var Registry = make(registry)

func (c registry) Add(name string, constructor Init) {
	c[name] = constructor
}
