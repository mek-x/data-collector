package sink

type Sink interface {
	Send(b []byte)
}

type SinkCfg struct {
	Name string
	Type string
	Spec string
}
