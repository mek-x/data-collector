package main

import (
	"bytes"
	"log"
	"os"
	"time"

	yaml "github.com/goccy/go-yaml"
	"gitlab.com/mek_x/DataCollector/pkg/datastore"
	"gitlab.com/mek_x/DataCollector/pkg/dispatcher/cron"
	"gitlab.com/mek_x/DataCollector/pkg/parser/jsonpath"
	"gitlab.com/mek_x/DataCollector/pkg/sink/stdout"
	"gitlab.com/mek_x/DataCollector/pkg/source/mqtt"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need config file!")
	}

	y, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("can't read config file: ", err)
	}

	ds := datastore.New()

	params := struct {
		Url  string `yaml:"url"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	}{}

	path, err := yaml.PathString("$.sources.mqtt.params")
	if err != nil {
		log.Fatal("path string is bad: ", err)
	}

	err = path.Read(bytes.NewReader(y), &params)
	if err != nil {
		log.Fatal("can't parse source params")
	}

	path, err = yaml.PathString("$.data.outside.path")
	if err != nil {
		log.Fatal("path string is bad: ", err)
	}

	var topic string
	err = path.Read(bytes.NewReader(y), &topic)
	if err != nil {
		log.Fatal("can't parse source params")
	}

	mqtt := mqtt.NewClient(params.Url, params.User, params.Pass)
	mqtt.Start()

	p := jsonpath.New("out", ds)
	p.AddVar("temp", "$.T")
	p.AddVar("humi", "$.H")
	p.AddVar("addr", "$.address")

	mqtt.AddDataSource(topic, p)

	d := cron.New("*/1 * * * * *", ds)
	d.AddSink("out", stdout.Stdout{})
	d.Start()

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
	}

	mqtt.End()
}
