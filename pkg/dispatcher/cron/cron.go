package cron

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"gitlab.com/mek_x/DataCollector/pkg/datastore"
	"gitlab.com/mek_x/DataCollector/pkg/dispatcher"
	"gitlab.com/mek_x/DataCollector/pkg/sink"
)

type cronDispatcher struct {
	cronString string
	sinks      map[string]sink.Sink
	ds         datastore.DataStore
}

var _ dispatcher.Dispatcher = (*cronDispatcher)(nil)

func New(cronString string, ds datastore.DataStore) *cronDispatcher {
	return &cronDispatcher{
		cronString: cronString,
		ds:         ds,
		sinks:      make(map[string]sink.Sink),
	}
}

func (c *cronDispatcher) sendToAll() {
	for val, s := range c.sinks {
		toSend, err := c.ds.Get(val)
		if err != nil {
			log.Println("can't find var in ds: ", err)
			continue
		}
		j, err := json.Marshal(toSend)
		if err != nil {
			log.Println("can't encode json: ", err)
			return
		}

		s.Send(j)
	}
}

func (c *cronDispatcher) Start() {
	s := gocron.NewScheduler(time.UTC)
	s.CronWithSeconds(c.cronString).Do(c.sendToAll)

	s.StartAsync()
}

func (c *cronDispatcher) AddSink(name string, s sink.Sink) {
	c.sinks[name] = s
}
