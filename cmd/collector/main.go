package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	yaml "github.com/goccy/go-yaml"
	"gitlab.com/mek_x/data-collector/pkg/collector"
	"gitlab.com/mek_x/data-collector/pkg/collector/mqtt"
	"gitlab.com/mek_x/data-collector/pkg/datastore"
	"gitlab.com/mek_x/data-collector/pkg/dispatcher/cron"
	"gitlab.com/mek_x/data-collector/pkg/parser/jsonpath"
	"gitlab.com/mek_x/data-collector/pkg/sink/stdout"
)

func parseConfig(in []byte, yamlPath string, object interface{}) error {

	path, err := yaml.PathString(yamlPath)
	if err != nil {
		return err
	}

	err = path.Read(bytes.NewReader(in), object)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need config file!")
	}

	y, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("can't read config file: ", err)
	}

	y = []byte(os.ExpandEnv(string(y)))

	ds := datastore.New()

	collectorsCfg := make(map[string]map[string]interface{})
	parseConfig(y, "$.collectors", &collectorsCfg)

	collectors := make(map[string]collector.Collector)

	for i, v := range collectorsCfg {
		switch v["type"].(string) {
		case "mqtt":
			params := mqtt.MqttParams{}
			parseConfig(y, fmt.Sprintf("$.collectors.%s.params", i), &params)
			mqtt := mqtt.NewClient(params)
			collectors[i] = mqtt
			log.Print("added collector: ", i, ", type: ", v["type"])
		default:
			log.Printf("config: collectors.%s - unknown type: %s", i, v["type"])
		}
	}

	var topic string
	err = parseConfig(y, "$.data.outside.path", &topic)
	if err != nil {
		log.Fatal("can't parse source params: ", err)
	}

	collectors["mqtt"].Start()

	p := jsonpath.New("out", ds)
	p.AddVar("temp", "$.T")
	p.AddVar("humi", "$.H")
	p.AddVar("addr", "$.address")

	collectors["mqtt"].AddDataSource(topic, p)

	d := cron.New("*/1 * * * * *", ds)
	d.AddSink("out", stdout.Stdout{})
	d.Start()

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
	}

	collectors["mqtt"].End()
}
