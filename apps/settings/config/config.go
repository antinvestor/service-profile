package config

import "github.com/pitabwire/frame"

type SettingsConfig struct {
	frame.ConfigurationDefault

	SecurelyRunService bool `default:"true" envconfig:"SECURELY_RUN_SERVICE"`
}
