package upgrade

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfg         *viper.Viper
	flagVersion string
)

var Command = &cobra.Command{
	Use:     "upgrade",
	Short:   "Quickly upgrade or cli version.",
	Long:    `Quickly upgrade or cli version.`,
	Example: `sudo listagent upgrade`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	initFlags()
}

func initFlags() {
	flgs := Command.Flags()
	flgs.StringVarP(&flagVersion, "version", "v", "", "upgrade (or downgrade) to the specified version.")
}
