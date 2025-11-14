package config

import "github.com/joeshaw/envdecode"

var StorageDirectory = "./storage/app/"

type Config struct {
	AppName                  string   `env:"APP_NAME"`
	AppVersion               string   `env:"APP_VERSION"`
	AppEnv                   string   `env:"APP_ENV,default=development"`
	ApiHost                  string   `env:"API_HOST"`
	ApiRpcPort               string   `env:"API_RPC_PORT"`
	ApiPort                  string   `env:"API_PORT,default=8760"`
	ApiDocPort               uint16   `env:"API_DOC_PORT,default=8761"`
	ShutdownTimeout          uint     `env:"API_SHUTDOWN_TIMEOUT_SECONDS,default=30"`
	AllowedCredentialOrigins []string `env:"ALLOWED_CREDENTIAL_ORIGINS"`
	MiddlewareAddress        string   `env:"MIDDLEWARE_ADDR"`
	JwtExpireDaysCount       int      `env:"JWT_EXPIRE_DAYS_COUNT"`
	MysqlOption
	RedisOption
}

// MysqlOption contains mySQL connection options
type MysqlOption struct {
	URI           string `env:"MYSQL_URI,default="`
	Pool          int    `env:"MYSQL_POOL,required"`
	SlowThreshold int    `env:"MYSQL_SLOW_LOG_THRESHOLD,required"`
}

// RedisOption contains Redis connection options
type RedisOption struct {
	Host     string `env:"REDIS_HOST,default=localhost"`
	Port     string `env:"REDIS_PORT,default=6379"`
	Password string `env:"REDIS_PASSWORD,default="`
	DB       int    `env:"REDIS_DB,default=0"`
}

	Uri          string `env:"MONGODB_URI,required"`
	DatabaseName string `env:"MONGODB_DATABASE_NAME,required"`
}

type RedisOption struct {
	Host           string `env:"REDIS_HOST,required"`
	Password       string `env:"REDIS_PASSWORD"`
	ReadTimeoutMs  int16  `env:"REDIS_READ_TIMEOUT,required"`
	WriteTimeoutMs int16  `env:"REDIS_WRITE_TIMEOUT,required"`
}

func NewConfig() *Config {
	var cfg Config
	if err := envdecode.Decode(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
