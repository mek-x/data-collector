package stdout

import (
	"fmt"

	"gitlab.com/mek_x/data-collector/pkg/sink"
)

type stdout struct{}

var _ sink.Sink = (*stdout)(nil)

func init() {
	sink.Registry.Add("stdout", New)
}

func New(params any) sink.Sink {
	return &stdout{}
}

func (*stdout) Send(b []byte) error {
	fmt.Println(string(b))
	return nil
}
