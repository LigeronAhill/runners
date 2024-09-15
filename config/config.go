package config

import "github.com/spf13/viper"

func Init(fileName string) (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigName(fileName)
	config.AddConfigPath(".")
	config.AddConfigPath("$HOME")
	err := config.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}
