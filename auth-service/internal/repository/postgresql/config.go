package pg

type Config struct {
	Port         int    `yaml:"port"`
	Host         string `yaml:"host"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"database"`
}
