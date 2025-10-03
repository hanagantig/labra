package config

import (
	"github.com/spf13/viper"
	"time"
)

type AppConfig struct {
	Name    string
	Env     string
	Version string
}

type HTTPConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type HTTPClientConfig struct {
	BaseURL string
	APIKey  string
}

type ThirdPartyService struct {
	Url string
}

type SQLConfig struct {
	Host              string
	Port              string
	DBName            string
	User              string
	Password          string
	Timeout           time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	InterpolateParams bool
	Charset           string
	ParseTime         bool
	Timezone          string
	Collation         string
	ConnMaxLifetime   time.Duration
	MaxOpenConns      int
	MaxIdleConns      int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	User     string
}

type AuthConfig struct {
	AccessTokenSecret  string
	AccessTokenTTL     time.Duration
	RefreshTokenSecret string
	RefreshTokenTTL    time.Duration
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SenderName   string
	SenderEmail  string
	AuthEmail    string
	AuthPassword string
}

type Config struct {
	App AppConfig

	Auth AuthConfig

	// Server
	HTTP HTTPConfig

	Email EmailConfig

	DocuPandaClient HTTPClientConfig

	MainDB SQLConfig

	Redis RedisConfig

	Service1 ThirdPartyService
}

func NewConfig(filePath string) (Config, error) {
	viper.SetConfigFile(filePath)

	conf := Config{}
	if err := viper.ReadInConfig(); err != nil {
		return conf, err
	}

	viper.SetDefault("app.name", "testapp")
	viper.SetDefault("app.env", "dev")
	viper.SetDefault("app.version", "v1")
	viper.SetDefault("http.read_timeout", "1s")
	viper.SetDefault("http.write_timeout", "1s")
	conf = Config{
		App: AppConfig{
			Name:    viper.GetString("app.name"),
			Env:     viper.GetString("app.env"),
			Version: viper.GetString("app.version"),
		},

		Auth: AuthConfig{
			AccessTokenSecret:  viper.GetString("auth.access_token_secret"),
			AccessTokenTTL:     viper.GetDuration("auth.access_token_ttl"),
			RefreshTokenSecret: viper.GetString("auth.refresh_token_secret"),
			RefreshTokenTTL:    viper.GetDuration("auth.refresh_token_ttl"),
		},

		HTTP: HTTPConfig{
			Host:         viper.GetString("http.host"),
			Port:         viper.GetString("http.port"),
			ReadTimeout:  viper.GetDuration("http.read_timeout"),
			WriteTimeout: viper.GetDuration("http.write_timeout"),
		},

		Email: EmailConfig{
			SMTPHost:     viper.GetString("email.smtp_host"),
			SMTPPort:     viper.GetInt("email.smtp_port"),
			SenderName:   viper.GetString("email.sender_name"),
			SenderEmail:  viper.GetString("email.sender_email"),
			AuthEmail:    viper.GetString("email.auth_email"),
			AuthPassword: viper.GetString("email.auth_password"),
		},

		DocuPandaClient: HTTPClientConfig{
			BaseURL: viper.GetString("docupanda_client.base_url"),
			APIKey:  viper.GetString("docupanda_client.api_key"),
		},

		Redis: RedisConfig{
			Host:     viper.GetString("redis.host"),
			Port:     viper.GetString("redis.port"),
			Password: viper.GetString("redis.password"),
			User:     viper.GetString("redis.user"),
		},

		MainDB: SQLConfig{
			Host:              viper.GetString("main_db.host"),
			Port:              viper.GetString("main_db.port"),
			DBName:            viper.GetString("main_db.name"),
			User:              viper.GetString("main_db.user"),
			Password:          viper.GetString("main_db.password"),
			ReadTimeout:       viper.GetDuration("main_db.read_timeout"),
			WriteTimeout:      viper.GetDuration("main_db.write_timeout"),
			Timeout:           viper.GetDuration("main_db.timeout"),
			InterpolateParams: false,
			Charset:           viper.GetString("main_db.charset"),
			ParseTime:         true,
			Timezone:          "UTC",
			Collation:         "",
		},
	}

	return conf, nil
}
