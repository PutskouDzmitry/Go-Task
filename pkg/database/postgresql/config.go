package postgresql

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"strconv"
	"task/internal/config"
	"time"
)

type postgresqlConfig struct {
	Host            string `yaml:"host"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Port            string `yaml:"port"`
	Database        string `yaml:"database"`
	MaxConn         int32  `yaml:"max_conn"`
	MaxIdleConn     string `yaml:"max_idle_conn"`
	MaxLifetimeConn string `yaml:"max_lifetime_conn"`
	MaxAttempts     string `yaml:"max_attempts"`
	MaxDelay        string `yaml:"max_delay"`
}

// newPostgresqlConfig initializes a new postgresqlConfig struct using the given config provider.
//
// The function takes a config.Provider as a parameter and returns a pointer to a postgresqlConfig struct and an error.
func newPostgresqlConfig(cfg *config.Config) (*postgresqlConfig, error) {
	return &postgresqlConfig{
		Host:            cfg.Postgresql.Host,
		Username:        cfg.Postgresql.Username,
		Password:        cfg.Postgresql.Password,
		Port:            cfg.Postgresql.Port,
		Database:        cfg.Postgresql.Database,
		MaxConn:         cfg.Postgresql.MaxConn,
		MaxIdleConn:     cfg.Postgresql.MaxIdleConn,
		MaxLifetimeConn: cfg.Postgresql.MaxLifetimeConn,
		MaxAttempts:     cfg.Postgresql.MaxAttempts,
		MaxDelay:        cfg.Postgresql.MaxDelay,
	}, nil
}

func (cfg *postgresqlConfig) getPostgresqlConfig(logger *zap.Logger) (*pgxpool.Config, error) {
	port, _ := strconv.ParseUint(cfg.Port, 10, 32)
	maxConnLifetime, _ := strconv.ParseUint(cfg.MaxLifetimeConn, 10, 8)
	maxIdleConn, _ := strconv.ParseUint(cfg.MaxIdleConn, 10, 8)
	healthCheckPeriod, _ := strconv.ParseUint(cfg.MaxAttempts, 10, 8)

	parseConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		logger.Error("failed to parse postgresql config, err: " + err.Error())
		return nil, err
	}

	parseConfig.ConnConfig.Host = cfg.Host
	parseConfig.ConnConfig.Port = uint16(port)
	parseConfig.ConnConfig.User = cfg.Username
	parseConfig.ConnConfig.Password = cfg.Password
	parseConfig.ConnConfig.Database = cfg.Database
	parseConfig.MaxConnLifetime = time.Duration(maxConnLifetime) * time.Second
	parseConfig.MaxConnIdleTime = time.Duration(maxIdleConn) * time.Second
	parseConfig.HealthCheckPeriod = time.Duration(healthCheckPeriod) * time.Second
	parseConfig.MaxConns = cfg.MaxConn

	return parseConfig, nil
}
