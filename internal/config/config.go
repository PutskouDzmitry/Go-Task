package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"sync"
	"time"
)

type Config struct {
	App struct {
		Name          string `yaml:"name" env-required:"true"`
		IsDevelopment bool   `yaml:"is-development"`
	} `yaml:"app"`
	Resources string `yaml:"resources"`
	HTTP      struct {
		Host         string        `yaml:"host" env-required:"true"`
		Port         int           `yaml:"port" env-required:"true"`
		ReadTimeout  time.Duration `yaml:"read-timeout" env-required:"true"`
		WriteTimeout time.Duration `yaml:"write-timeout" env-required:"true"`
	}
	Postgresql struct {
		Host            string `yaml:"host" env-required:"true"`
		Username        string `yaml:"username" env-required:"true"`
		Password        string `yaml:"password" env-required:"true"`
		Port            string `yaml:"port" env-required:"true"`
		Database        string `yaml:"database" env-required:"true"`
		MaxConn         int32  `yaml:"max_conn" env-required:"true"`
		MaxIdleConn     string `yaml:"max_idle_conn" env-required:"true"`
		MaxLifetimeConn string `yaml:"max_lifetime_conn" env-required:"true"`
		MaxAttempts     string `yaml:"max_attempts" env-required:"true"`
		MaxDelay        string `yaml:"max_delay" env-required:"true"`
	} `yaml:"postgresql"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}

		if err := cleanenv.ReadConfig(os.Getenv("CONFIG_PATH"), instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Print(help)
			log.Fatal(err)
		}
	})

	return instance
}
