package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"
)

const (
	// appName is the application name, which will also be used as the binary name.
	appName = "shockrays"

	// relPath is the default relative path to the application's home directory e.g. ~/.shockrays.
	relPath = ".shockrays"

	// cfgName is the default name of the configuration file.
	cfgName = "config.yaml"
)

// defaultHomeDir is the default home directory for the application.
// It is initialized with an OS specific value when the main function is called and the Execute call is invoked.
var defaultHomeDir string

// NewRootCmd returns the root command for the application. Subcommands are added to the root command,
// along with any persistent flags that should be exposed across the application.
func NewRootCmd(log *slog.Logger, homeDir string) *cobra.Command {
	defaultHomeDir = filepath.Join(homeDir, relPath)
	defaultOutputDir = filepath.Join(defaultHomeDir, defaultOutputDir)
	defaultPRDir = filepath.Join(defaultHomeDir, defaultPRDir)

	if runtime.GOOS == "windows" {
		defaultPRBinName = defaultPRBinName + ".exe"
	}

	a := &appState{
		log:      log,
		viper:    viper.New(),
		homePath: defaultHomeDir,
	}

	var rootCmd = &cobra.Command{
		Use:   appName,
		Short: "shockrays is a tool for interacting with ProjectorRays.",
		Long:  "shockrays is a tool for interacting with, it offers some improvements when dumping scripts from a directory of Director files.",
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		// this takes effect after flags are parsed.
		// which ensures log based flags are consumed when the logger is initialized.
		if log == nil {
			logLvl := new(slog.LevelVar)
			logLvl.Set(slog.LevelInfo)

			debug := a.viper.GetBool(flagDebug)
			if debug {
				logLvl.Set(slog.LevelDebug)
			}

			a.log = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLvl}))
		}

		return a.loadConfigFile()
	}

	rootCmd = homeFlag(a.viper, rootCmd, a.homePath)
	rootCmd = debugFlag(a.viper, rootCmd, a.debug)

	rootCmd.AddCommand(
		configCmd(a),
		decompileCmd(a),
	)

	return rootCmd
}

// Execute is called by the main function, it is used to register the subcommands.
// This function should only be called once.
func Execute(homeDir string) {
	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd(nil, homeDir)
	rootCmd.SilenceUsage = true

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		// Wait for interrupt signal.
		sig := <-sigCh

		// Cancel context on root command.
		// If the invoked command respects this quickly, the main goroutine will quit right away.
		cancel()

		// Short delay before printing the received signal message.
		// This should result in cleaner output from non-interactive commands that stop quickly.
		time.Sleep(250 * time.Millisecond)
		fmt.Fprintf(os.Stderr, "Received signal %v. Attempting clean shutdown. Send interrupt again to force hard shutdown.\n", sig)

		// Dump all goroutines on panic, not just the current one.
		debug.SetTraceback("all")

		// Block waiting for a second interrupt or a timeout.
		// The main goroutine ought to finish before either case is reached.
		// But if a case is reached, panic so that we get a non-zero exit and a dump of remaining goroutines.
		select {
		case <-time.After(time.Minute):
			panic(errors.New("shockrays did not shut down within one minute of interrupt"))
		case sig := <-sigCh:
			panic(fmt.Errorf("received signal %v; forcing termination of shockrays", sig))
		}
	}()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

// withUsage wraps a PositionalArgs to display usage only when the PositionalArgs
// variant is violated.
func withUsage(inner cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := inner(cmd, args); err != nil {
			cmd.Root().SilenceUsage = false
			cmd.SilenceUsage = false
			return err
		}

		return nil
	}
}
