package iotplotter

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"gitlab.com/mek_x/data-collector/pkg/sink"
)

const baseUrl = "http://iotplotter.com/api/v2/feed"

type IotPlotterParams struct {
	Apikey string
	Feed   string
}

type iotplotter struct {
	apikey string
	feed   string
}

var _ sink.Sink = (*iotplotter)(nil)

func init() {
	sink.Registry.Add("iotplotter", New)
}

func New(p any) sink.Sink {

	var opt IotPlotterParams

	if err := mapstructure.Decode(p, &opt); err != nil {
		return nil
	}

	if opt.Apikey == "" || opt.Feed == "" {
		return nil
	}

	return &iotplotter{
		apikey: opt.Apikey,
		feed:   opt.Feed,
	}
}

func (p *iotplotter) Send(b []byte) error {
	r, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", baseUrl, p.feed), bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("api-key", p.apikey)

	client := &http.Client{}
	_, err = client.Do(r)

	return err
}
