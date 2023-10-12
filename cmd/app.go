package cmd

import (
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path"
)

// appState represents the application state that can be mutated during configuration and startup.
type appState struct {
	log      *slog.Logger
	viper    *viper.Viper
	homePath string
	debug    bool
	config   Config
}

// loadConfigFile attempts to read the configuration file from disk and initialize an instance of Config.
func (a *appState) loadConfigFile() error {
	p := path.Join(a.homePath, "config", cfgName)

	if _, err := os.Stat(p); err != nil {
		// This func is called in PersistentPreRunE on the root cmd,
		// which means the function is called before Run for each cmd.
		// We need to return nil here because it's possible the config init cmd has not been executed.
		return nil
	}

	f, err := os.ReadFile(p)
	if err != nil {
		return err
	}

	cfg := Config{}
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return err
	}

	a.config = cfg

	return nil
}
