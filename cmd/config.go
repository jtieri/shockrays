package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strings"
)

var (
	// defaultOutputDir is the default directory where Lingo scripts will be dumped.
	defaultOutputDir = "output"

	// defaultPRDir is the default directory where the ProjectorRays binary should be located.
	defaultPRDir = "projector-rays"

	// defaultPRBinName is the default name of the ProjectorRays binary.
	defaultPRBinName = "projectorrays"
)

// Config represents the configurable options available for the application.
type Config struct {
	// PRDirectory is the full path name to the directory where the ProjectorRays binary should be located.
	PRDirectory string `yaml:"projectorrays-dir" json:"projectorrays-dir"`

	// ScriptOutputDir is the full path name to the directory where scripts should be dumped.
	// Inside the directory the path points to, scripts will be dumped in a subdirectory named after the file
	// that was decompiled.
	ScriptOutputDir string `yaml:"script-output-dir" json:"script-output-dir"`

	// PRBinName is the name of the ProjectorRays binary.
	PRBinName string `yaml:"projectorrays-bin-name" json:"projectorrays-bin-name"`
}

// MustYAML returns Config serialized as a YAML string. This method is used where an error cannot be recovered from,
// so if an error is encountered a panic is invoked.
func (c *Config) MustYAML() []byte {
	out, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	return out
}

// configCmd registers the subcommand tree for the various configuration based commands.
func configCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "Subcommands for managing the application's config file",
	}

	cmd.AddCommand(
		configInitCmd(a),
		configShowCmd(a),
	)
	return cmd
}

// configInitCmd is a command that initializes the default config file at the location specified by --home.
func configInitCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Initializes a default config file",
		Long:    "Initializes a default config file at the location specified by the --home flag",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s config init --home %s
$ %s cfg i`, appName, a.homePath, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				return err
			}

			cfgDir := path.Join(home, "config")
			cfgPath := path.Join(cfgDir, "config.yaml")

			// Check if the config file already exists in the default config directory.
			// If not, check if the config directory already exists.
			// If not, check if the home directory already exists.
			// If not, make the home directory and then make the config directory.
			if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
				if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
					if _, err := os.Stat(home); os.IsNotExist(err) {
						if err = os.Mkdir(home, os.ModePerm); err != nil {
							return err
						}
					}

					if err = os.Mkdir(cfgDir, os.ModePerm); err != nil {
						return err
					}
				}

				f, err := os.Create(cfgPath)
				if err != nil {
					return err
				}
				defer f.Close()

				defaultCfg := defaultConfig()

				if _, err = f.Write(defaultCfg.MustYAML()); err != nil {
					return err
				}

				// Create the default projectorrays directory if it does not exist.
				if _, err := os.Stat(defaultCfg.PRDirectory); os.IsNotExist(err) {
					if err = os.Mkdir(defaultCfg.PRDirectory, os.ModePerm); err != nil {
						return err
					}
				}

				// Create the default script output directory if it does not exist.
				if _, err := os.Stat(defaultCfg.ScriptOutputDir); os.IsNotExist(err) {
					if err = os.Mkdir(defaultCfg.ScriptOutputDir, os.ModePerm); err != nil {
						return err
					}
				}

				return nil
			}

			return fmt.Errorf("config file already exists: %s", cfgPath)
		},
	}

	return cmd
}

// configShowCmd is a command that attempts to read the config file off disk and prints its contents to the console.
func configShowCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"s"},
		Short:   "Prints current configuration",
		Args:    withUsage(cobra.NoArgs),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s config show --home %s
$ %s cfg list`, appName, a.homePath, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				return err
			}

			cfgPath := path.Join(home, "config", "config.yaml")

			if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
				if _, err := os.Stat(home); os.IsNotExist(err) {
					return fmt.Errorf("path specified by --home flag does not exist: %s", home)
				}

				return fmt.Errorf("config file does not exist: %s", cfgPath)
			}

			out, err := yaml.Marshal(a.config)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(out))
			return nil
		},
	}

	return cmd
}

// defaultConfig returns the default Config for the application.
func defaultConfig() Config {
	return Config{
		PRDirectory:     defaultPRDir,
		ScriptOutputDir: defaultOutputDir,
		PRBinName:       defaultPRBinName,
	}
}
