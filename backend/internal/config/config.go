package config

import "booky-backend/internal/utils"

type DatabaseConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

type Config struct {
	DBCfg  *DatabaseConfig
	SvPort string
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
		SvPort: utils.GetEnvOrDefault("PORT", ":8080"),
	}
}
