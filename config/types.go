package config

type ConfigObject struct {
	MinioRoot MinioRootConfig
	Network   NetworkConfig
	Logger    LogConfig
}

type MinioRootConfig struct {
	User string
	Pass string
}

type NetworkConfig struct {
	Host           string
	Port           string
	MaxConnections int
	TimeoutSeconds int
}

type LogConfig struct {
	FilePath string
	Level    int
}
