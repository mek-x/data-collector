package gotify

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/mitchellh/mapstructure"
	"gitlab.com/mek_x/data-collector/pkg/sink"
)

type GotifyParams struct {
	Url      string
	Token    string
	Title    string
	Priority int
}

type gotify struct {
	url      string
	token    string
	title    string
	priority int
}

var _ sink.Sink = (*gotify)(nil)

func init() {
	sink.Registry.Add("gotify", New)
}

func New(p any) sink.Sink {

	var opt GotifyParams

	if err := mapstructure.Decode(p, &opt); err != nil {
		return nil
	}

	if opt.Url == "" || opt.Token == "" {
		return nil
	}

	if opt.Title == "" {
		opt.Title = ">>><<<"
	}

	return &gotify{
		url:      opt.Url,
		token:    opt.Token,
		title:    opt.Title,
		priority: opt.Priority,
	}
}

func (g *gotify) Send(b []byte) error {

	fMap := template.FuncMap{
		"now": func(f string) string { return time.Now().Format(f) },
	}

	tmpl, err := template.New("title").Funcs(fMap).Parse(g.title)
	if err != nil {
		return err
	}

	var buf string
	title := bytes.NewBufferString(buf)

	// Run the template to verify the output.
	err = tmpl.Execute(title, nil)
	if err != nil {
		return err
	}

	_, err = http.PostForm(fmt.Sprintf("%s/message?token=%s", g.url, g.token),
		url.Values{
			"message":  {string(b)},
			"title":    {title.String()},
			"priority": {fmt.Sprint(g.priority)},
		})

	return err
}
