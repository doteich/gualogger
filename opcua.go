package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

var (
	last_keepalive time.Time
	retry_count    = 10
)

func (o *OpcConfig) InitSuperVisor(ctx context.Context) {

	_, err := o.Connection.CreateClient(ctx)

	if err != nil {
		Logger.Error(err.Error(), "func", "InitSuperVisor")
		return
	}

	Logger.Info(fmt.Sprintf("successfully connected to opcua on endpoint %s:%d", o.Connection.Endpoint, o.Connection.Port))

	// for {

	// }

}

func (c *OpcConnection) CreateClient(ctx context.Context) (*opcua.Client, error) {
	eps, err := opcua.GetEndpoints(ctx, fmt.Sprintf("opc.tcp://%s:%d", c.Endpoint, c.Port))

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

	client, err := opcua.NewClient(c.Endpoint, opts...)

	if err != nil {
		return nil, err
	}

	return client, nil

}
