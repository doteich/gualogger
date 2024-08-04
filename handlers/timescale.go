package handlers

import "fmt"

type TimeScaleDB struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

var ()

func (t *TimeScaleDB) CreateConPool() error {
	fmt.Println(t)
	return nil
}
