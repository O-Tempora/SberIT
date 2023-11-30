package config

type Config struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	DbHost string `yaml:"dbhost"`
	DbPort int    `yaml:"dbport"`
	DbBase string `yaml:"dbbase"`
}
