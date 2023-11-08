package dispatcher

import "gitlab.com/mek_x/data-collector/pkg/sink"

type Dispatcher interface {
	Start()
	AddSink(name string, sink sink.Sink)
}
