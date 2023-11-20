package sink

type Sink interface {
	Send(b []byte) error
}

type SinkCfg struct {
	Name string
	Type string
	Spec string
}

type Init func(params any) Sink
type registry map[string]Init

var Registry = make(registry)

func (s registry) Add(name string, constructor Init) {
	s[name] = constructor
}
