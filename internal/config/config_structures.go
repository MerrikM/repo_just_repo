package config

type ServerConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	BasePath string `yaml:"base_path"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

type JWTConfig struct {
	SecretKey       string `yaml:"secret_key"`
	AccessTokenTTL  string `yaml:"access_token_ttl"`
	RefreshTokenTTL string `yaml:"refresh_token_ttl"`
}

type WebhookConfig struct {
	URL     string `yaml:"url"`
	Timeout string `yaml:"timeout"`
}

type AdminConfig struct {
	AdminToken string `yaml:"admin_token"`
}

type TTL struct {
	S3AndRedis int `yaml:"s3_and_redis"`
}

type MINIO struct {
	Login    string `yaml:"minio_login"`
	Password string `yaml:"minio_password"`
}
