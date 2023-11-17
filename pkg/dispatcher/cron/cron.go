package cron

import (
	"bytes"
	"encoding/json"
	"log"
	"text/template"
	"time"

	"github.com/go-co-op/gocron"
	"gitlab.com/mek_x/data-collector/pkg/datastore"
	"gitlab.com/mek_x/data-collector/pkg/dispatcher"
	"gitlab.com/mek_x/data-collector/pkg/sink"

	"github.com/antonmedv/expr"
)

type sinkInstance struct {
	iface sink.Sink
	cfg   sink.SinkCfg
}

type cronDispatcher struct {
	cronString string
	sinks      []sinkInstance
	ds         datastore.DataStore
}

var _ dispatcher.Dispatcher = (*cronDispatcher)(nil)

func New(cronString string, ds datastore.DataStore) *cronDispatcher {
	return &cronDispatcher{
		cronString: cronString,
		ds:         ds,
		sinks:      make([]sinkInstance, 0),
	}
}

func (c *cronDispatcher) sendToAll() {
	for _, s := range c.sinks {

		var toSend []byte

		switch s.cfg.Type {
		case "expr":
			env := c.ds.Map()
			program, err := expr.Compile(s.cfg.Spec, expr.Env(env))
			if err != nil {
				log.Println("can't compile expr: ", err)
				continue
			}

			output, err := expr.Run(program, env)
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
		s.iface.Send(toSend)
	}
}

func (c *cronDispatcher) Start() {
	s := gocron.NewScheduler(time.UTC)
	s.CronWithSeconds(c.cronString).Do(c.sendToAll)

	s.StartAsync()
}

func (c *cronDispatcher) AddSink(s sink.Sink, cfg sink.SinkCfg) {
	sink := sinkInstance{
		iface: s,
		cfg:   cfg,
	}
	c.sinks = append(c.sinks, sink)
}
