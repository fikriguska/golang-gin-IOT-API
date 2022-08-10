package config

type Configuration struct {
	Port   int    `env:"PORT", required`
	DBHost string `env:"DB_HOST", required`
	DBUser string `env:"DB_USER", required`
	DBPass string `env:"DB_PASS", required`
	DBName string `env:"DB_NAME", required`
	DBPort int    `env:"DB_PORT", required`
}
