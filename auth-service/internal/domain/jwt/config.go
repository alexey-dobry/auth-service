package jwt

import "time"

type Config struct {
	AccessSecret  string `validate:"required" yaml:"access-secret"`
	RefreshSecret string `validate:"required" yaml:"refresh-secret"`
	TTL           TTL    `validate:"required" yaml:"ttl"`
}

type TTL struct {
	AccessTTL  time.Duration `validate:"required" yaml:"access-ttl"`
	RefreshTTL time.Duration `validate:"required" yaml:"refresh-ttl"`
}
