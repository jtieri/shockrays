package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagHome  = "home"
	flagDebug = "debug"
	flagPRDir = "projector-rays"
)

// homeFlag defines the flag for the home directory of the application and binds the default path for the directory.
func homeFlag(v *viper.Viper, cmd *cobra.Command, homePath string) *cobra.Command {
	cmd.PersistentFlags().StringVar(&homePath, flagHome, defaultHomeDir, "set home directory")
	if err := v.BindPFlag(flagHome, cmd.PersistentFlags().Lookup(flagHome)); err != nil {
		panic(err)
	}
	return cmd
}

// debugFlag defines the flag for enabling debug logging in the application and binds the default value to false.
func debugFlag(v *viper.Viper, cmd *cobra.Command, debug bool) *cobra.Command {
	cmd.PersistentFlags().BoolVarP(&debug, flagDebug, "d", false, "enable debug output")
	if err := v.BindPFlag("debug", cmd.PersistentFlags().Lookup("debug")); err != nil {
		panic(err)
	}
	return cmd
}

// projectorRaysFlag defines the flag for specifying the directory where the ProjectorRays binary should be located.
func projectorRaysFlag(v *viper.Viper, cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringP(flagPRDir, "p", defaultPRDir, "default directory for ProjectorRays binary")
	if err := v.BindPFlag(flagPRDir, cmd.Flags().Lookup(flagPRDir)); err != nil {
		panic(err)
	}
	return cmd
}
