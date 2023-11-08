package jsonpath

import (
	"encoding/json"
	"log"

	"gitlab.com/mek_x/DataCollector/pkg/datastore"
	"gitlab.com/mek_x/DataCollector/pkg/parser"

	jp "github.com/PaesslerAG/jsonpath"
)

type jsonpathParser struct {
	store datastore.DataStore
	vars  map[string]string
	name  string
}

var _ parser.Parser = (*jsonpathParser)(nil)

func New(name string, store datastore.DataStore) *jsonpathParser {
	return &jsonpathParser{
		name:  name,
		store: store,
		vars:  make(map[string]string),
	}
}

func (j *jsonpathParser) AddVar(name, v string) {
	j.vars[name] = v
}

func (j *jsonpathParser) Parse(buf []byte) error {
	v := interface{}(nil)

	err := json.Unmarshal(buf, &v)
	if err != nil {
		log.Println("bad json: ", err)
		return err
	}

	out := make(map[string]interface{})

	for k, val := range j.vars {
		a, err := jp.Get(val, v)
		if err != nil {
			log.Println("bad jsonpath: ", err)
			goto fail
		}
		out[k] = a
	}

	j.store.Publish(j.name, out)

fail:
	return err
}
