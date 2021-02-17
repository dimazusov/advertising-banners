package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	cfg, err := New("config_test.yaml")
	require.Nil(t, err)
	require.NotNil(t, cfg)

	require.Equal(t, "127.0.0.1", cfg.Server.HTTP.Host)
	require.Equal(t, "8080", cfg.Server.HTTP.Port)

	require.Equal(t, "debug", cfg.Logger.Level)
	require.Equal(t, "/log/log.txt", cfg.Logger.Path)

	require.Equal(t, "postgres", cfg.DB.Postgres.Dialect)
	require.Equal(t, "host=postgres port=5432 dbname=db user=db password=db sslmode=disable", cfg.DB.Postgres.Dsn)

	require.Equal(t, uint(20), cfg.DB.Memory.MaxSize)

	require.Equal(t, "kafka:9092", cfg.Kafka.Address)
	require.Equal(t, "30ms", cfg.Kafka.BatchTimeout)
	require.Equal(t, "30s", cfg.Kafka.DialerTimeout)
	require.Equal(t, "main_group", cfg.Kafka.Group)

	require.Equal(t, "postgres", cfg.Repository.Type)
}
