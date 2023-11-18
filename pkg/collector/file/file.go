package file

import (
	"os"
	"time"

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

var _ collector.Collector = (*fileSource)(nil)

func New(interval int) *fileSource {
	return &fileSource{
		sources:  make([]source, 0),
		interval: interval,
		end:      make(chan bool),
	}
}

func (f *fileSource) Start() error {
	go func() {
		for {
			select {
			case <-f.end:
				f.end <- false
				close(f.end)
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
