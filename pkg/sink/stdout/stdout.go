package stdout

import (
	"fmt"

	"gitlab.com/mek_x/data-collector/pkg/sink"
)

type Stdout struct{}

var _ sink.Sink = (*Stdout)(nil)

func (Stdout) Send(b []byte) {
	fmt.Println(string(b))
}
