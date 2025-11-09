package handlers

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gualogger/logging"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

// Kafka struct holds the configuration for the Kafka client
type Redpanda struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	Auth    struct {
		SASL struct {
			Type string `mapstructure:"type"`
			User string `mapstructure:"user"`
			Pass string `mapstructure:"pass"`
		} `mapstructure:"sasl"`
	} `mapstructure:"auth"`
	TLS struct {
		InsecureSkipVerify bool `mapstructure:"insecure_skip_verify"`
	} `mapstructure:"tls"`
	Client *kgo.Client
}

// NewKafkaClient creates a new Kafka client with the given configuration
func (r *Redpanda) Initialize(ctx context.Context) error {

	opts := []kgo.Opt{
		kgo.SeedBrokers(r.Brokers...),
		kgo.WithLogger(logging.NewKgoLogger(logging.Logger)),
	}

	if r.TLS.InsecureSkipVerify {
		tlsCfg := new(tls.Config)
		tlsCfg.InsecureSkipVerify = true
		opts = append(opts, kgo.DialTLSConfig(tlsCfg))
	} else {
		opts = append(opts, kgo.DialTLS())
	}

	if r.Auth.SASL.Type != "" {
		switch r.Auth.SASL.Type {
		case "scram-sha-256":
			opts = append(opts, kgo.SASL(scram.Auth{
				User: r.Auth.SASL.User,
				Pass: r.Auth.SASL.Pass,
			}.AsSha256Mechanism()))
		case "scram-sha-512":
			opts = append(opts, kgo.SASL(scram.Auth{
				User: r.Auth.SASL.User,
				Pass: r.Auth.SASL.Pass,
			}.AsSha512Mechanism()))
		case "plain":
			opts = append(opts, kgo.SASL(plain.Auth{
				User: r.Auth.SASL.User,
				Pass: r.Auth.SASL.Pass,
			}.AsMechanism()))
		}
	}

	client, err := kgo.NewClient(opts...)

	if err != nil {
		return err
	}

	if err := client.Ping(ctx); err != nil {
		return err
	}

	r.Client = client
	return nil
}

func (r *Redpanda) Publish(ctx context.Context, p Payload) error {
	// 1. Add a timeout to the context
	produceCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	b, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	rec := &kgo.Record{
		Key:       []byte(p.Id),
		Topic:     r.Topic,
		Timestamp: p.TS,
		Value:     b,
	}

	results := r.Client.ProduceSync(produceCtx, rec)

	// 2. Correctly check for errors
	if err := results.FirstErr(); err != nil {
		return fmt.Errorf("produce sync failed: %v", err)
	}

	return nil
}

func (r *Redpanda) Shutdown(ctx context.Context) error {
	r.Client.Close()
	return nil
}

func (r *Redpanda) Ping(ctx context.Context) error {

	ctx_t, done := context.WithTimeout(ctx, 10*time.Second)

	defer done()

	if err := r.Client.Ping(ctx_t); err != nil {
		return err
	}

	return nil
}
