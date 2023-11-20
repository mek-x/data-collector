package mqtt

import (
	"crypto/tls"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mitchellh/mapstructure"
	"gitlab.com/mek_x/data-collector/internal/utils"
	"gitlab.com/mek_x/data-collector/pkg/collector"
	"gitlab.com/mek_x/data-collector/pkg/parser"
)

type mqttSource struct {
	client      mqtt.Client
	mqttOptions *mqtt.ClientOptions
}

type MqttParams struct {
	Url  string
	User string
	Pass string
}

var _ collector.Collector = (*mqttSource)(nil)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("MQTT received message: %s from topic: %s", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("MQTT connected to broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("MQTT connection lost: %v", err)
}

func init() {
	collector.Registry.Add("mqtt", New)
}

func New(p any) collector.Collector {

	var opt MqttParams

	if err := mapstructure.Decode(p, &opt); err != nil {
		return nil
	}

	ssl := tls.Config{
		RootCAs: nil,
	}

	var m mqttSource

	opts := mqtt.NewClientOptions()
	opts.AddBroker(opt.Url)
	opts.SetClientID("data-collector-" + utils.RandomString(5))
	opts.SetUsername(opt.User)
	opts.SetPassword(opt.Pass)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetTLSConfig(&ssl)

	client := mqtt.NewClient(opts)

	m.mqttOptions = opts
	m.client = client

	return &m
}

func (m *mqttSource) connect() error {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *mqttSource) End() {
	m.client.Disconnect(100)
}

func (m *mqttSource) Start() error {
	return m.connect()
}

func (m *mqttSource) AddDataSource(topic string, parser parser.Parser) error {

	log.Printf("adding data source: %s", topic)
	if token := m.client.Subscribe(topic, 0, func(c mqtt.Client, msg mqtt.Message) {
		if err := parser.Parse(msg.Payload()); err != nil {
			log.Println("can't parse: ", err)
		}
	}); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
