package sink

type Sink interface {
	Send(b []byte)
}
