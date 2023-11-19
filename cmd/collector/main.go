package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	yaml "github.com/goccy/go-yaml"
	"gitlab.com/mek_x/data-collector/pkg/collector"
	"gitlab.com/mek_x/data-collector/pkg/collector/file"
	"gitlab.com/mek_x/data-collector/pkg/collector/mqtt"
	"gitlab.com/mek_x/data-collector/pkg/datastore"
	"gitlab.com/mek_x/data-collector/pkg/dispatcher"
	"gitlab.com/mek_x/data-collector/pkg/dispatcher/cron"
	"gitlab.com/mek_x/data-collector/pkg/parser"
	"gitlab.com/mek_x/data-collector/pkg/parser/jsonpath"
	"gitlab.com/mek_x/data-collector/pkg/sink"
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

	var version int
	parseConfig(y, "$.version", &version)
	if version != 1 {
		log.Fatal("config file has not supported version: ", version)
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
			collectors[i] = mqtt.NewClient(params)
			log.Print("added collector: ", i, ", type: ", v["type"])
		case "file":
			var interval int
			parseConfig(y, fmt.Sprintf("$.collectors.%s.params", i), &interval)
			collectors[i] = file.New(interval)
			log.Print("added collector: ", i, ", type: ", v["type"])
		default:
			log.Printf("config: collectors.%s - unknown type: %s", i, v["type"])
		}
	}

	sinksCfg := make(map[string]map[string]interface{})
	parseConfig(y, "$.sinks", &sinksCfg)
	sinks := make(map[string]sink.Sink)

	for i, v := range sinksCfg {
		switch v["type"].(string) {
		case "stdout":
			s := stdout.Stdout{}
			sinks[i] = s
			log.Print("added sink: ", i, ", type: ", v["type"])
		default:
			log.Printf("config: sinks.%s - unknown type: %s", i, v["type"])
		}
	}

	dispatchersCfg := make([]map[string]interface{}, 0)
	parseConfig(y, "$.dispatchers", &dispatchersCfg)
	dispatchers := make([]dispatcher.Dispatcher, 0)

	for i, v := range dispatchersCfg {
		switch v["type"].(string) {
		case "cron":
			cronStr := v["param"].(string)
			c := cron.New(cronStr, ds)
			targets := make([]sink.SinkCfg, 0)
			parseConfig(y, fmt.Sprintf("$.dispatchers[%d].sinks", i), &targets)
			for _, s := range targets {
				c.AddSink(sinks[s.Name], s)
				log.Print("cron - added sink: ", s)
			}

			dispatchers = append(dispatchers, c)
			log.Print("added dispatcher: ", i, ", type: ", v["type"])
		case "event":
		default:
			log.Printf("config: dispatchers[%d] - unknown type: %s", i, v["type"])
		}
	}

	for _, c := range collectors {
		c.Start()
	}

	dataCfg := make(map[string]map[string]interface{})
	parseConfig(y, "$.data", &dataCfg)

	for i, v := range dataCfg {
		var p parser.Parser
		switch v["parser"].(string) {
		case "jsonpath":
			p = jsonpath.New(i, ds)
			vars := make(map[string]string)
			parseConfig(y, fmt.Sprintf("$.data.%s.vars", i), &vars)
			for k, val := range vars {
				p.AddVar(k, val)
			}
			conv := make(map[string]string)
			parseConfig(y, fmt.Sprintf("$.data.%s.conv", i), &conv)
			for k, val := range conv {
				p.AddConv(k, val)
			}
		default:
			log.Printf("config: data[%s] - unknown parser type: %s", i, v["parser"])
			continue
		}

		topic := v["path"].(string)
		c := v["collector"].(string)

		collectors[c].AddDataSource(topic, p)
	}

	for _, d := range dispatchers {
		d.Start()
	}

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
	}

	for _, c := range collectors {
		c.End()
	}
}
