package config

import "github.com/pitabwire/frame/v3/config"

type SettingsConfig struct {
	config.ConfigurationDefault

	SecurelyRunService bool `default:"true" envconfig:"SECURELY_RUN_SERVICE"`
}
