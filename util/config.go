package util

import "github.com/spf13/viper"

type Config struct {
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBSource   string `mapstructure:"DB_SOURCE"`
	ServerAddr string `mapstructure:"SERVER_ADDR"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
