package cmd

import (
	"fmt"
	"github.com/rad12000/list-agent/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "manage and view the current cli version",
	Long:  "manage and view the current cli version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Current: %s\nLatest: %s\n", version.Version(), version.Latest())
	},
}
