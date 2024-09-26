package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
    Postgres PostgresConfig
    Redis    RedisConfig
    Server   ServerConfig
}

type PostgresConfig struct {
    URL             string
    MaxIdleConns    int
    MaxOpenConns    int
    ConnMaxLifetime int
}

type RedisConfig struct {
    Addr                  string
    Pass                  string
    DB                    int
    CacheExpirationMinutes int
}

type ServerConfig struct {
    Port        int
    ReadTimeout int
}

var Cfg Config

func InitConfig() {
    viper.SetConfigName("config")
    viper.SetConfigType("json")
    viper.AddConfigPath("./config")

    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("Error reading config file: %s", err)
    }

    // 读取 PostgreSQL 配置
    Cfg.Postgres.URL = viper.GetString("postgres.url")
    Cfg.Postgres.MaxIdleConns = viper.GetInt("postgres.max_idle_conns")
    Cfg.Postgres.MaxOpenConns = viper.GetInt("postgres.max_open_conns")
    Cfg.Postgres.ConnMaxLifetime = viper.GetInt("postgres.conn_max_lifetime")

    // 读取 Redis 配置
    Cfg.Redis.Addr = viper.GetString("redis.addr")
    Cfg.Redis.Pass = viper.GetString("redis.pass")
    Cfg.Redis.DB = viper.GetInt("redis.db")
    Cfg.Redis.CacheExpirationMinutes = viper.GetInt("redis.cache_expiration_minutes")

    // 读取服务器配置
    Cfg.Server.Port = viper.GetInt("server.port")
    Cfg.Server.ReadTimeout = viper.GetInt("server.read_timeout")
}
