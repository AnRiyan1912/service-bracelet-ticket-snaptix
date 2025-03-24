package config

type Config struct {
	Host             string   `env:"HOST"`
	Port             int      `env:"PORT" envDefault:"3000"`
	IsDevelopment    bool     `env:"IS_DEVELOPMENT"`
	ProxyHeader      string   `env:"PROXY_HEADER"`
	LogFields        []string `env:"LOG_FIELDS" envSeparator:","`
	DbDsn            string   `env:"DB_DSN"`
	DbDriver         string   `env:"DB_DRIVER"`
	EncriptionKey    string   `env:"ENCRYPTION_KEY"`
	RedisClusterMode bool     `env:"REDIS_CLUSTER_MODE"`
	RedisAddress     []string `env:"REDIS_ADDR"`
	RedisPassword    string   `env:"REDIS_PASSWORD"`
	RedisUrl         string   `env:"REDIS_URL"`
}
