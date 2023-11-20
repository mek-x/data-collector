package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	yaml "github.com/goccy/go-yaml"
	"gitlab.com/mek_x/data-collector/pkg/collector"
	"gitlab.com/mek_x/data-collector/pkg/datastore"
	"gitlab.com/mek_x/data-collector/pkg/dispatcher"
	"gitlab.com/mek_x/data-collector/pkg/parser"
	"gitlab.com/mek_x/data-collector/pkg/sink"

	_ "gitlab.com/mek_x/data-collector/internal/modules"
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
		collectorType := v["type"].(string)
		var params interface{}
		err := parseConfig(y, fmt.Sprintf("$.collectors.%s.params", i), &params)
		if errors.Is(err, yaml.ErrNotFoundNode) {
			params = nil
		} else if err != nil {
			log.Print(i, ": error parsing collector params: ", err)
			continue
		}
		collectorInit, ok := collector.Registry[collectorType]
		if !ok {
			log.Print(i, ": unknown collector type - ", collectorType)
			continue
		}
		c := collectorInit(params)
		if c == nil {
			log.Print(i, ": can't initialize collector, config error?")
			continue
		}
		collectors[i] = c
		log.Print("added collector: ", i, ", type: ", collectorType)
	}

	sinksCfg := make(map[string]map[string]interface{})
	parseConfig(y, "$.sinks", &sinksCfg)
	sinks := make(map[string]sink.Sink)

	for i, v := range sinksCfg {
		sinkType := v["type"].(string)
		sinkInit, ok := sink.Registry[sinkType]
		if !ok {
			log.Print(i, ": unknown sink type - ", sinkType)
			continue
		}
		s := sinkInit()
		if s == nil {
			log.Print(i, ": can't initialize sink, config error?")
		}
		sinks[i] = s
		log.Print("added sink: ", i, ", type: ", sinkType)
	}

	dispatchersCfg := make([]map[string]interface{}, 0)
	parseConfig(y, "$.dispatchers", &dispatchersCfg)
	dispatchers := make([]dispatcher.Dispatcher, 0)

	for i, v := range dispatchersCfg {
		dispatcherType := v["type"].(string)
		param, ok := v["param"]
		if !ok {
			param = nil
		}

		dispatcherInit, ok := dispatcher.Registry[dispatcherType]
		if !ok {
			log.Print(i, ": unknown dispatcher type - ", dispatcherType)
			continue
		}

		d := dispatcherInit(param, ds)
		if d == nil {
			log.Print(i, ": can't initialize dispatcher, config error?")
		}

		targets := make([]sink.SinkCfg, 0)
		parseConfig(y, fmt.Sprintf("$.dispatchers[%d].sinks", i), &targets)
		for _, s := range targets {
			d.AddSink(sinks[s.Name], s)
			log.Print(i, ": added sink - ", s)
		}

		dispatchers = append(dispatchers, d)
		log.Print("added dispatcher: ", i, ", type: ", dispatcherType)
	}

	for _, c := range collectors {
		c.Start()
	}

	dataCfg := make(map[string]map[string]interface{})
	parseConfig(y, "$.data", &dataCfg)

	for i, v := range dataCfg {
		parserType := v["parser"].(string)

		parserInit, ok := parser.Registry[parserType]
		if !ok {
			log.Print(i, ": unknown parser type - ", parserType)
			continue
		}
		p := parserInit(i, ds)
		if p == nil {
			log.Print(i, ": can't initialize parser, config error?")
		}

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

		path := v["path"].(string)
		c := v["collector"].(string)

		collectors[c].AddDataSource(path, p)
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
