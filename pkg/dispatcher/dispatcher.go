package dispatcher

import "gitlab.com/mek_x/DataCollector/pkg/sink"

type Dispatcher interface {
	Start()
	AddSink(name string, sink sink.Sink)
}
