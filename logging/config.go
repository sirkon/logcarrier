package logging

// Config настройки логирования
type Config struct {
	Logfile string `yaml:"logfile"`
	Level   string `yaml:"level"` // default:"debug"
}

// NewConfig возвращает инстанс Config
func NewConfig() *Config {
	return &Config{
		Level: "info",
	}
}
