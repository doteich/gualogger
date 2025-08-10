package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	Paho "github.com/eclipse/paho.mqtt.golang"
)

type Mqtt struct {
	Protocol        string `mapstructure:"protocol"`
	ProtocolVersion uint   `mapstructure:"protocol_version"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Topic           string `mapstructure:"topic"`
	ClientID        string `mapstructure:"client_id"`
	QoS             int    `mapstructure:"qos"`
	Retain          bool   `mapstructure:"retain"`
	TopicByNodeId   bool   `mapstructure:"topic_by_nodeid"`
	Client          Paho.Client
}

func (m *Mqtt) Initialize(ctx context.Context, cb func(context.Context) []Payload) error {
	opts := Paho.NewClientOptions()

	opts.ClientID = m.ClientID

	if m.Protocol != "mqtt" && m.Protocol != "mqtts" {
		return fmt.Errorf("invalid protocol: %s", m.Protocol)
	}

	opts.AddBroker(fmt.Sprintf("%s://%s:%d", m.Protocol, m.Host, m.Port))
	if m.Username != "" || m.Password != "" {
		opts.SetUsername(m.Username)
		opts.SetPassword(m.Password)
	}
	opts.SetProtocolVersion(m.ProtocolVersion)
	opts.AutoReconnect = true
	opts.SetKeepAlive(60)

	c := Paho.NewClient(opts)

	c.Connect()
	m.Client = c

	return nil
}

func (m *Mqtt) Publish(ctx context.Context, p Payload) error {
	topic := m.Topic

	if m.TopicByNodeId && len(p.Meta) > 0 {
		topic = p.Meta[0].Value
	}

	b, err := json.Marshal(p)

	if err != nil {
		return fmt.Errorf("failed to marshal payload: %s", err)
	}

	if token := m.Client.Publish(topic, byte(m.QoS), m.Retain, b); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish message: %s", token.Error())
	}
	return nil
}
func (m *Mqtt) Shutdown(ctx context.Context) error {
	m.Client.Disconnect(100)
	return nil
}
