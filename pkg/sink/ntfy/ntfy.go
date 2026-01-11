package ntfy

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/mitchellh/mapstructure"
	"gitlab.com/mek_x/data-collector/pkg/sink"
)

type NtfyParams struct {
	Url      string
	Topic    string
	Token    string
	Title    string
	Priority int
}

type ntfy struct {
	url      string
	topic    string
	token    string
	title    string
	priority int
}

var _ sink.Sink = (*ntfy)(nil)

func init() {
	sink.Registry.Add("ntfy", New)
}

func New(p any) sink.Sink {

	var opt NtfyParams

	if err := mapstructure.Decode(p, &opt); err != nil {
		return nil
	}

	if opt.Url == "" || opt.Topic == "" {
		return nil
	}

	if opt.Priority == 0 {
		opt.Priority = 3
	}

	return &ntfy{
		url:      opt.Url,
		topic:    opt.Topic,
		token:    opt.Token,
		title:    opt.Title,
		priority: opt.Priority,
	}
}

func reformat(in string) (string, error) {
	fMap := template.FuncMap{
		"now": func(f string) string { return time.Now().Format(f) },
	}

	tmpl, err := template.New("title").Funcs(fMap).Parse(in)
	if err != nil {
		return "", err
	}

	var buf string
	title := bytes.NewBufferString(buf)

	// Run the template to verify the output.
	err = tmpl.Execute(title, nil)
	if err != nil {
		return "", err
	}

	return title.String(), nil
}

func (n *ntfy) Send(b []byte) error {

	r, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", n.url, n.topic), bytes.NewBuffer(b))

	if err != nil {
		return err
	}

	if n.token != "" {
		r.Header.Add("Authorization", "Bearer "+n.token)
	}

	if n.title != "" {
		title, err := reformat(n.title)
		if err != nil {
			return err
		}

		r.Header.Add("Title", title)
	}

	if n.priority != 0 {
		r.Header.Add("Priority", strconv.Itoa(n.priority))
	}

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed, code: %d", resp.StatusCode)
	}

	return nil
}
