package jwt

import "time"

type Config struct {
	AccessSecret  string `yaml:"access-secret"`
	RefreshSecret string `yaml:"refresh-secret"`
	TTL           TTL    `yaml:"ttl"`
}

type TTL struct {
	AccessTTL  time.Duration `yaml:"access-ttl"`
	RefreshTTL time.Duration `yaml:"refresh-ttl"`
}
