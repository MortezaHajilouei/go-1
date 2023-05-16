package config

// Service details
type Service struct {
	Name    string `yaml:"usecase.name" required:"true"`
	ID      uint32 `yaml:"usecase.id" required:"true"`
	BaseURL string `yaml:"usecase.baseURL"`
	GRPC    struct {
		Host     string `yaml:"grpc.host"`
		Port     string `yaml:"grpc.port"`
		TLS      bool   `yaml:"grpc.tls"`
		Protocol string `yaml:"protocol"`
	}
	HTTP struct {
		Host           string `yaml:"http.host"`
		Port           string `yaml:"http.port"`
		RequestTimeout string `yaml:"http.requestTimeout"`
	}
	Router []Router `yaml:"usecase.routers"`
}

type Router struct {
	Description       string   `yaml:"description"`
	Method            string   `yaml:"method"`
	MaxAllowedAnomaly float32  `yaml:"maxAllowedAnomaly"`
	Middlewares       []string `yaml:"middlewares"`
}
