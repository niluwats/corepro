package utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbDriver      string `mapstructure:"DB_DRIVER"`
	DBUsername    string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	ServerPort    string `mapstructure:"SERVER_PORT"`
	EmailFrom     string `mapstructure:"EMAIL_FROM"`
	EmailPassword string `mapstructure:"EMAIL_PASSWORD"`
	EmailPort     string `mapstructure:"EMAIL_PORT"`
	JwtSecretKey  string `mapstructure:"JWT_SECRET_KEY"`
	SmsDigest     string `mapstructure:"SMS_DIGEST"`
	SmsUser       string `mapstructure:"SMS_USER"`
	SmsAuthKey    string `mapstructure:"SMS_AUTH_KEY"`
	Mask          string `mapstructure:"MASK"`
	CampaignName  string `mapstructure:"CAMPAIGN_NAME"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
