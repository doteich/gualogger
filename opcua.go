package main

import (
	"context"
	"fmt"
	"gualogger/handlers"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/monitor"
	"github.com/gopcua/opcua/ua"
)

var (
	last_keepalive      time.Time
	con_active          bool
	retry_count         int
	current_retry_count = 0
	Subs                map[uint32]*monitor.Subscription
)

func (o *OpcConfig) InitSuperVisor(ctx context.Context) {

	retry_count = o.Connection.Retries

	Subs = make(map[uint32]*monitor.Subscription)

	c, err := o.Connection.CreateClient(ctx)

	if err != nil {
		Logger.Error(err.Error(), "func", "InitSuperVisor")
		return
	}

	Logger.Info(fmt.Sprintf("successfully connected to opcua on endpoint %s:%d", o.Connection.Endpoint, o.Connection.Port))

	subctx, cancel := context.WithCancel(ctx)

	if err := InitSubs(c, ctx, subctx, &o.Subscription.Nodeids, o.Subscription.Interval); err != nil {
		Logger.Error(fmt.Sprintf("error while creating node monitor: %s", err.Error()), "func", "InitSuperVisor")
		return
	}

	con_active = true

	for {

		time.Sleep(3 * time.Duration(o.Subscription.Interval))

		if time.Since(last_keepalive) > time.Duration(6*o.Subscription.Interval)*time.Second {

			current_retry_count++

			if retry_count < current_retry_count {

				Logger.Warn(fmt.Sprintf("maximum number of %d retries exceeded- shutting down", retry_count), "func", "InitSuperVisor")
				cancel()
				ctx.Done()
				break
			}

			con_active = false

			Logger.Warn(fmt.Sprintf("received last keepalive over %d seconds ago attempting retry attempt %d/%d", 60, current_retry_count, retry_count), "func", "InitSuperVisor")

			if con_active {
				cancel()
				c.Close(ctx)
			}

			c, err = o.Connection.CreateClient(ctx)

			if err != nil {
				Logger.Error(err.Error(), "func", "InitSuperVisor")
				continue
			}

			subctx, cancel = context.WithCancel(ctx)

			if err := InitSubs(c, ctx, subctx, &o.Subscription.Nodeids, o.Subscription.Interval); err != nil {
				Logger.Error(fmt.Sprintf("error while creating node monitor: %s", err.Error()), "func", "InitSuperVisor")
				continue
			}
			Logger.Info("connection retry successful")
			con_active = true
			current_retry_count = 0
		}

	}

}

func (c *OpcConnection) CreateClient(ctx context.Context) (*opcua.Client, error) {

	con_string := fmt.Sprintf("opc.tcp://%s:%d", c.Endpoint, c.Port)

	eps, err := opcua.GetEndpoints(ctx, con_string)

	if err != nil {
		return nil, err
	}

	if len(eps) < 1 {
		return nil, fmt.Errorf("no endpoints found - check configuration")
	}

	ep := opcua.SelectEndpoint(eps, c.Policy, ua.MessageSecurityModeFromString(c.Mode))

	opts := []opcua.Option{
		opcua.ApplicationName("guanaco"),
		opcua.AutoReconnect(true),
		opcua.ReconnectInterval(10 * time.Second),
		opcua.SecurityPolicy(c.Policy),
		opcua.SecurityMode(ua.MessageSecurityModeFromString(c.Mode)),
	}

	switch c.Authentication.Type {

	case "User&Password":
		opts = append(opts, opcua.AuthUsername(c.Authentication.Credentials.Username, c.Authentication.Credentials.Password))
		opts = append(opts, opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeUserName))

	case "Certificate": // Cert Config goes here, but has to be evaluated first

	default:
		opts = append(opts, opcua.AuthAnonymous())
		opts = append(opts, opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeAnonymous))
	}

	if c.Policy != "None" {
		if c.Certificate.AutoCreate {
			if err := CreateKeyPair(); err != nil {
				return nil, err
			}
		}

		opts = append(opts, opcua.CertificateFile("./certs/cert.pem"))
		opts = append(opts, opcua.PrivateKeyFile("./certs/key.pem"))
	}

	client, err := opcua.NewClient(con_string, opts...)

	if err != nil {
		return nil, err
	}

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	return client, nil

}

func InitSubs(c *opcua.Client, pctx context.Context, ctx context.Context, ids *[]string, iv int) error {
	m, err := monitor.NewNodeMonitor(c)

	if err != nil {

		return err
	}

	go CreateSubscription(pctx, ctx, m, ids, iv)

	time.Sleep(60 * time.Second)
	return nil
}

func CreateSubscription(pctx context.Context, ctx context.Context, m *monitor.NodeMonitor, ids *[]string, iv int) {

	sub, err := m.Subscribe(pctx, &opcua.SubscriptionParameters{Interval: time.Duration(iv) * time.Second},
		func(s *monitor.Subscription, dcm *monitor.DataChangeMessage) {
			if dcm.Error != nil {
				Logger.Error(fmt.Sprintf("error with received sub message: %s - nodeid %s", dcm.Error.Error(), dcm.NodeID))
			} else if dcm.Status != ua.StatusOK {
				Logger.Error(fmt.Sprintf("received bad status for sub message: %s - nodeid %s", dcm.Status, dcm.NodeID))
			} else {
				var dt string
				switch dcm.Value.Value().(type) {
				case int:
					dt = "Int"
				case uint8:
					dt = "u8"
				case uint16:
					dt = "u16"
				case uint32:
					dt = "u32"
				case int8:
					dt = "i8"
				case int16:
					dt = "i16"
				case int32:
					dt = "i32"
				case int64:
					dt = "i64"
				case float32:
					dt = "f32"
				case float64:
					dt = "f64"
				case bool:
					dt = "Bool"
				default:
					dt = "Str"

				}
				if dcm.NodeID.String() == "i=2258" {
					last_keepalive = time.Now()
				} else {
					p := handlers.Payload{Value: dcm.Value.Value(), TS: dcm.SourceTimestamp, Name: dcm.NodeID.StringID(), Id: dcm.NodeID.String(), Datatype: dt}

					mgr.Publish(ctx, p)

				}

			}

		})

	if err != nil {
		Logger.Error(fmt.Sprintf("error while creating subscription: %s", err.Error()))
		return
	}

	for _, n := range *ids {
		_, err := sub.AddMonitorItems(ctx, monitor.Request{NodeID: ua.MustParseNodeID(n), MonitoringMode: ua.MonitoringModeReporting, MonitoringParameters: &ua.MonitoringParameters{DiscardOldest: true, QueueSize: 1}})
		if err != nil {
			Logger.Error(fmt.Sprintf("error adding subscription item: %s", err.Error()))
			continue
		}
	}

	_, err = sub.AddMonitorItems(ctx, monitor.Request{NodeID: ua.MustParseNodeID("i=2258"), MonitoringMode: ua.MonitoringModeReporting, MonitoringParameters: &ua.MonitoringParameters{DiscardOldest: true, QueueSize: 1}})

	if err != nil {
		Logger.Error(fmt.Sprintf("error adding subscription item: %s", err.Error()))
		return
	}

	id := sub.SubscriptionID()
	Subs[id] = sub

	Logger.Info(fmt.Sprintf("successfully initialized subscription with id:%d", id))

	defer TerminateSub(pctx, sub, id)
	<-ctx.Done()
}

func TerminateSub(ctx context.Context, s *monitor.Subscription, id uint32) {

	Logger.Warn(fmt.Sprintf("terminating subscription with id: %d - delivered: %d - dropped: %d", id, s.Delivered(), s.Dropped()))
	delete(Subs, id)
	s.Unsubscribe(ctx)

}
