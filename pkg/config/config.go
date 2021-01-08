// pkg/config/config.go

package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
)

var (
	Config = config{}
)

type config struct {
	DBUser     string `mapstructure:"DB_USER"`
	DBPass     string `mapstructure:"DB_PASS"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
	ApiPort    string `mapstructure:"API_PORT"`
	AppAccessSecret  string `mapstructure:"APP_ACCESS_SECRET"`
	AppRefreshSecret  string `mapstructure:"APP_REFRESH_SECRET"`
	Migrate    string
}

func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&Config)

	flag.StringVar(&Config.Migrate,
		"migrate", "up",
		"specify if we should be migrating DBConn 'up' or 'down'")

	flag.Parse()
	return
}

func (c *config) GetDBConnStr() string {
	return c.getDBConnStr(c.DBHost, c.DBName)
}

func (c *config) getDBConnStr(dbhost, dbname string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser,
		c.DBPass,
		dbhost,
		c.DBPort,
		dbname,
	)
}

func (c *config) GetApiPort() string {
	return c.ApiPort
}

func (c *config) GetMigration() string {
	return c.Migrate
}
