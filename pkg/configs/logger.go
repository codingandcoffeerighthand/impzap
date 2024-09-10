package configs

type LogConfig struct {
	LogLevel    string `yaml:"log_level"`
	Dir         string `yaml:"dir"`
	MaxBackups  int    `yaml:"max_backups"`
	MaxSize     int    `yaml:"max_size"`
	MaxAge      int    `yaml:"max_age"`
	Compress    bool   `yaml:"compress"`
	ShowConsole bool   `yaml:"show_console"`
}
