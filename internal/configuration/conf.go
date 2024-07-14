package configuration

import "github.com/spf13/viper"

func LoadConfig() {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("./config")
}
