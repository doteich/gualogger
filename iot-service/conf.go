package main

import (
	"gualogger/handlers"

	"github.com/spf13/viper"
)

type Configuration struct {
	Opcua    OpcConfig         `mapstructure:"opcua"`
	Redpanda handlers.Redpanda `mapstructure:"redpanda"`
}

type OpcConfig struct {
	Connection   OpcConnection `mapstructure:"connection"`
	Subscription Subscription  `mapstructure:"subscription"`
}

type Subscription struct {
	Nodeids  []Nodeid `mapstructure:"nodeids"`
	Interval int      `mapstructure:"sub_interval"`
}

type OpcConnection struct {
	Endpoint       string            `mapstructure:"endpoint"`
	Port           int               `mapstructure:"port"`
	Mode           string            `mapstructure:"mode"`
	Policy         string            `mapstructure:"policy"`
	Authentication OpcAuthentication `mapstructure:"authentication"`
	Certificate    OpcCerts          `mapstructure:"certificate"`
	Retries        int               `mapstructure:"retry_count"`
}

type OpcAuthentication struct {
	Type        string `mapstructure:"type"`
	Credentials struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"credentials"`
	Certificate struct {
		CertificatePath string `mapstructure:"certificate_path"`
		PrivateKeyPath  string `mapstructure:"private_key_path"`
	} `mapstructure:"certificate"`
}

type OpcCerts struct {
	AutoCreate      bool   `mapstructure:"auto_create"`
	CertificatePath string `mapstructure:"certificate_path"`
	PrivateKeyPath  string `mapstructure:"private_key_path"`
}

type Nodeid struct {
	Id   string          `mapstructure:"id"`
	Meta []handlers.Meta `mapstructure:"meta"`
}

func LoadConfig() (*Configuration, error) {

	var conf Configuration

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/config")   // Linux FS
	v.AddConfigPath("$HOME/.config") // Windows FS
	v.AddConfigPath("./configs")     // Local Testing

	if err := v.ReadInConfig(); err != nil {
		return &conf, err
	}

	if err := v.Unmarshal(&conf); err != nil {
		return &conf, err
	}

	return &conf, nil
}

// Returns a map of all possible Exporters
// To add a new Exporter add a new entry in format [`conf key name`]=Exporter struct
