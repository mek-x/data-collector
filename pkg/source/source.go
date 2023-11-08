package source

import (
	"gitlab.com/mek_x/DataCollector/pkg/parser"
)

type Source interface {
	Start() error
	AddDataSource(path string, parser parser.Parser) error
}
