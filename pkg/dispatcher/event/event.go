package event

import (
	"bytes"
	"encoding/json"
	"log"
	"text/template"
	"time"

	"gitlab.com/mek_x/data-collector/pkg/datastore"
	"gitlab.com/mek_x/data-collector/pkg/dispatcher"
	"gitlab.com/mek_x/data-collector/pkg/sink"

	"github.com/antonmedv/expr"
	"github.com/mitchellh/mapstructure"
)

type sinkInstance struct {
	iface sink.Sink
	cfg   sink.SinkCfg
}

type EventParams struct {
	Trigger string
	Var     string
	Expr    string
}

type eventDispatcher struct {
	eventParams EventParams
	sinks       []sinkInstance
	ds          datastore.DataStore
}

var _ dispatcher.Dispatcher = (*eventDispatcher)(nil)

func init() {
	dispatcher.Registry.Add("event", New)
}

func New(param any, ds datastore.DataStore) dispatcher.Dispatcher {
	var opt EventParams

	if err := mapstructure.Decode(param, &opt); err != nil {
		return nil
	}

	return &eventDispatcher{
		eventParams: opt,
		ds:          ds,
		sinks:       make([]sinkInstance, 0),
	}
}

type triggerExpr struct {
	new any
	old any
}

func (c *eventDispatcher) sendToAll(key string, t time.Time, v, old interface{}) {

	for _, s := range c.sinks {

		var toSend []byte

		switch s.cfg.Type {
		case "expr":
			output, err := expr.Eval(s.cfg.Spec, c.ds.Map())
			if err != nil {
				log.Println("can't run expr: ", err)
				continue
			}

			toSend, err = json.Marshal(output)
			if err != nil {
				log.Println("can't encode json: ", err)
				continue
			}
		case "template":
			t := template.Must(template.New("msg").Parse(s.cfg.Spec))
			b := &bytes.Buffer{}
			err := t.Execute(b, c.ds.Map())
			if err != nil {
				log.Print("can't execute template: ", err)
				continue
			}
			toSend = b.Bytes()
		default:
			log.Printf("unknown sink (%s) data type: %s", s.cfg.Name, s.cfg.Type)
			continue
		}

		log.Print("Sending to: ", s.cfg.Name)
		err := s.iface.Send(toSend)
		if err != nil {
			log.Print("Sink error: ", err)
		}
	}
}

func (c *eventDispatcher) Start() {

}

func (c *eventDispatcher) AddSink(s sink.Sink, cfg sink.SinkCfg) {
	sink := sinkInstance{
		iface: s,
		cfg:   cfg,
	}
	c.sinks = append(c.sinks, sink)
	c.ds.Register([]string{c.eventParams.Trigger}, c.sendToAll)
}
