package windy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"gitlab.com/mek_x/data-collector/pkg/sink"
)

const baseUrl = "https://stations.windy.com/pws/update"

type WindyParams struct {
	Apikey string
	Id     int
}

type windy struct {
	apikey string
	id     int
}

type observation struct {
	Temp      float64 `json:"temp"`
	Humi      float64 `json:"rh"`
	Pressure  float64 `json:"mbar"`
	Timestamp int     `json:"ts"`
	Station   int     `json:"station"`
}

type messageFormat struct {
	Observations []observation `json:"observations"`
}

var _ sink.Sink = (*windy)(nil)

func init() {
	sink.Registry.Add("windy", New)
}

func New(p any) sink.Sink {

	var opt WindyParams

	if err := mapstructure.Decode(p, &opt); err != nil {
		return nil
	}

	if opt.Apikey == "" || opt.Id == 0 {
		return nil
	}

	return &windy{
		apikey: opt.Apikey,
		id:     opt.Id,
	}
}

func (w *windy) Send(b []byte) error {
	var data observation
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	msg := messageFormat{
		Observations: []observation{{
			Temp:      data.Temp,
			Humi:      data.Humi,
			Pressure:  data.Pressure,
			Timestamp: data.Timestamp,
			Station:   w.id,
		}},
	}

	out, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	r, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", baseUrl, w.apikey), bytes.NewBuffer(out))

	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(r)

	return err
}
