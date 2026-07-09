package config

import "github.com/pitabwire/frame/v2/config"

type SettingsConfig struct {
	config.ConfigurationDefault

	SecurelyRunService bool `default:"true" envconfig:"SECURELY_RUN_SERVICE"`
}
