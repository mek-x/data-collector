package source

import (
	"gitlab.com/mek_x/DataCollector/pkg/datastore"
	"gitlab.com/mek_x/DataCollector/pkg/parser"
)

type Source interface {
	Start(store *datastore.DataStore) error
	AddDataSource(path string, parser parser.Parser) error
}
