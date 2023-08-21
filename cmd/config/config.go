package config

import "github.com/spf13/viper"

type Config struct {
	HttpPort string   `mapstructure:"HTTP_PORT"`
	Database Database `mapstructure:",squash"`
	Redis    Redis    `mapstructure:",squash"`
}

type Database struct {
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	Name     string `mapstructure:"DB_NAME"`
}

type Redis struct {
	URI      string `mapstructure:"REDIS_URI"`
	Password string `mapstructure:"REDIS_PASSWORD"`
}

func Load() (Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.ReadInConfig()

	viper.BindEnv("HTTP_PORT")
	viper.BindEnv("DB_USERNAME")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("REDIS_URI")
	viper.BindEnv("REDIS_PASSWORD")

	var c Config
	err := viper.Unmarshal(&c)

	return c, err
}
