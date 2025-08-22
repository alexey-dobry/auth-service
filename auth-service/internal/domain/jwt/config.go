package jwt

import (
	"crypto/ecdsa"
	"time"
)

type Config struct {
	AccessSecret  *ecdsa.PrivateKey `validate:"required" yaml:"access-secret"`
	RefreshSecret *ecdsa.PrivateKey `validate:"required" yaml:"refresh-secret"`
	TTL           TTL               `validate:"required" yaml:"ttl"`
}

type TTL struct {
	AccessTTL  time.Duration `validate:"required" yaml:"access-ttl"`
	RefreshTTL time.Duration `validate:"required" yaml:"refresh-ttl"`
}
