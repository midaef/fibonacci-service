package config

// Config ...
type Config struct {
	AppConfig AppConfig
	Redis     Redis
}

// AppConfig ...
type AppConfig struct {
	IP       string `yaml:"ip"`
	GRPCPort string `yaml:"grpc_port"`
	HTTPPort string `yaml:"http_port"`
}

// Redis ...
type Redis struct {
	DB       int    `yaml:"db"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
}
