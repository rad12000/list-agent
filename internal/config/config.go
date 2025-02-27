package config

import (
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func Directory() string {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	return filepath.Join(home, ".listagent")
}
