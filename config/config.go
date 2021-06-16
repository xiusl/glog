package config

import "github.com/spf13/viper"

type GLogConfig struct {
	EtcdAddr  string `mapstructure:"ETCD_ADDR"`
	KafkaAddr string `mapstructure:"KAFKA_ADDR"`
	EsAddr    string `mapstructure:"ES_ADDR"`
	EtcdKey   string `mapstructure:"ETCD_KEY"`
}

func LoadConfig() (config GLogConfig, err error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
