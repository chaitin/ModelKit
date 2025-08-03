package config

import (
	_ "embed"
	"strings"

	"github.com/spf13/viper"

	"github.com/chaitin/ModelKit/backend/pkg/logger"
)

//go:embed config.json.tmpl
var ConfigTmpl []byte

type Config struct {
	Debug bool `mapstructure:"debug"`

	Logger *logger.Config `mapstructure:"logger"`

	Database struct {
		Master          string `mapstructure:"master"`
		Slave           string `mapstructure:"slave"`
		MaxOpenConns    int    `mapstructure:"max_open_conns"`
		MaxIdleConns    int    `mapstructure:"max_idle_conns"`
		ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	} `mapstructure:"database"`

}

func Init() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("MODELKIT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("debug", false)
	v.SetDefault("logger.level", "info")
	v.SetDefault("database.master", "")
	v.SetDefault("database.slave", "")
	v.SetDefault("database.max_open_conns", 50)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", 30)

	c := Config{}
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
