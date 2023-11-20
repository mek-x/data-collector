package file

import (
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
	"gitlab.com/mek_x/data-collector/pkg/collector"
	"gitlab.com/mek_x/data-collector/pkg/parser"
)

type source struct {
	path   string
	parser parser.Parser
}

type fileSource struct {
	sources  []source
	interval int
	end      chan bool
}

/* Main config structure. Descibes `params` field from configuration yaml. */
type FileParams struct {
	Interval int
}

var _ collector.Collector = (*fileSource)(nil)

/* In this module init function we register our collector in global collector registry */
func init() {
	collector.Registry.Add("file", New)
}

func New(p any) collector.Collector {
	var opt FileParams

	if err := mapstructure.Decode(p, &opt); err != nil {
		return nil
	}

	// Set defaults
	if opt.Interval == 0 {
		opt.Interval = 60
	}

	return &fileSource{
		sources:  make([]source, 0),
		interval: opt.Interval,
		end:      make(chan bool),
	}
}

func (f *fileSource) Start() error {
	go func() {
	loop:
		for {
			select {
			case <-f.end:
				f.end <- false
				close(f.end)
				break loop
			case <-time.After(time.Duration(f.interval) * time.Second):
				for _, i := range f.sources {
					buf, err := os.ReadFile(i.path)
					if err != nil {
						continue
					}
					i.parser.Parse(buf)
				}
			}
		}
	}()
	return nil
}

func (f *fileSource) AddDataSource(path string, parser parser.Parser) error {
	f.sources = append(f.sources, source{
		path:   path,
		parser: parser,
	})
	return nil
}

func (f *fileSource) End() {
	f.end <- true
	<-f.end
}
