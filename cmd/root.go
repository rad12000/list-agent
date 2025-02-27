package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logLevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "listagent",
	Short: "A cli tool that makes searching through listing agents, such as zillow, easier",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var level slog.Level
		switch strings.ToLower(logLevel) {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelInfo
		}

		fmt.Println("using log level", level)
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   false,
			Level:       level,
			ReplaceAttr: nil,
		})))
	},
	Long: ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "The minimum log level. Default so info. Valid values are: debug, info, warn, and error.")
	rootCmd.AddCommand(zillowCmd, versionCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(filepath.Join(home, ".listagent"))
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.ConfigFileUsed()
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
