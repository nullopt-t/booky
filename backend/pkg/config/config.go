package config

import "booky-backend/pkg/utils"

type DatabaseConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type Config struct {
	DBCfg        *DatabaseConfig
	RedisCfg     *RedisConfig
	SvPort       string
	JwtSecretKey string
}

func Load() *Config {
	return &Config{
		DBCfg: &DatabaseConfig{
			DBHost:     utils.GetEnvOrDefault("DB_HOST", "localhost"),
			DBPort:     utils.GetEnvOrDefault("DB_PORT", "5432"),
			DBUser:     utils.GetEnvOrDefault("DB_USER", "bookshop"),
			DBPassword: utils.GetEnvOrDefault("DB_PASSWORD", "bookshop123"),
			DBName:     utils.GetEnvOrDefault("DB_NAME", "bookshop"),
		},
		RedisCfg: &RedisConfig{
			Addr:     utils.GetEnvOrDefault("REDIS_ADDR", "localhost:6379"),
			Password: utils.GetEnvOrDefault("REDIS_PASSWORD", ""),
			DB:       utils.GetEnvAsInt("REDIS_DB", 0),
		},
		SvPort:       utils.GetEnvOrDefault("PORT", "8080"),
		JwtSecretKey: utils.GetEnvOrDefault("JWT_SECRET", "jwt-secret-key"),
	}
}
