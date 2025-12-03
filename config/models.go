package config

type Config struct {
	Server   ServerConfig   `yaml:"Server" mapstructure:"Server"`
	Postgres PostgresConfig `yaml:"Postgres" mapstructure:"Postgres"`
	Monitor  MonitorConfig  `yaml:"Monitor" mapstructure:"Monitor"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
	DBName   string `yaml:"DBName" mapstructure:"DBName"`
	SSLMode  string `yaml:"sslMode" mapstructure:"sslMode"`
	PgDriver string `yaml:"pgDriver" mapstructure:"pgDriver"`
}

type ServerConfig struct {
	AppVersion string `yaml:"appVersion" mapstructure:"appVersion"`
	Host       string `yaml:"host" mapstructure:"host" validate:"required"`
	Port       string `yaml:"port" mapstructure:"port" validate:"required"`
}

type MonitorConfig struct {
	IntervalSeconds       int `yaml:"intervalSeconds"`
	RequestTimeoutSeconds int `yaml:"requestTimeoutSeconds"`
	HistoryLimit          int `yaml:"historyLimit"`
	MaxConcurrency        int `yaml:"maxConcurrency"`
}
