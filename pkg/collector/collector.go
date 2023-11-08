package collector

import (
	"gitlab.com/mek_x/data-collector/pkg/parser"
)

type Collector interface {
	Start() error
	AddDataSource(path string, parser parser.Parser) error
}
