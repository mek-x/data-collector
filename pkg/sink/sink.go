package sink

type Sink interface {
	Send(b []byte)
}

type SinkCfg struct {
	Name string
	Type string
	Spec string
}

type Init func() Sink
type sinkRegistry map[string]Init

var Registry = make(sinkRegistry)

func (s sinkRegistry) Add(name string, constructor Init) {
	s[name] = constructor
}
