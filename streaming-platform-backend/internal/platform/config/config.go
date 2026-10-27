package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

func Load[T any]() (T, error) {
	var config T
	if err := envconfig.Process("", &config); err != nil {
		var zero T
		return zero, fmt.Errorf("load config: %w", err)
	}
	return config, nil
}
