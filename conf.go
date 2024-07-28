package main

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Opcua     OpcConfig `mapstructure:"opcua"`
	Exporters Exporters `mapstructure:"exporters"`
}

type OpcConfig struct {
	Connection   OpcConnection `mapstructure:"connection"`
	Subscription Subscription  `mapstructure:"subscription"`
}

type Subscription struct {
	Nodeids  []string `mapstructure:"nodeids"`
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

type Exporters struct {
}

func LoadConfig() (Configuration, error) {

	var conf Configuration

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/gopclogs")   // Linux FS
	v.AddConfigPath("$HOME/.gopclogs") // Windows FS
	v.AddConfigPath("./configs")       // Local Testing

	if err := v.ReadInConfig(); err != nil {
		return conf, err
	}

	if err := v.Unmarshal(&conf); err != nil {
		return conf, err
	}

	return conf, nil
}
