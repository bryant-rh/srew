package cmd

import (
	"fmt"
	"os"

	"github.com/bryant-rh/srew/cmd/client/global"
	"github.com/bryant-rh/srew/pkg/installation"
	"github.com/bryant-rh/srew/pkg/installation/receipt"
	"github.com/bryant-rh/srew/pkg/installation/validation"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCmdUpgrade() *cobra.Command {
	upgradeCmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade installed plugins to newer versions",
		Long: `Upgrade installed plugins to a newer version.
This will reinstall all plugins that have a newer version in the local index.
To only upgrade single plugins provide them as arguments:
srew upgrade foo bar"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var ignoreUpgraded bool
			var skipErrors bool

			var pluginNames []string
			if len(args) == 0 {
				// Upgrade all plugins.
				installed, err := installation.GetInstalledPluginReceipts(paths.InstallReceiptsPath())
				if err != nil {
					return errors.Wrap(err, "failed to find all installed versions")
				}
				for _, receipt := range installed {
					pluginNames = append(pluginNames, receipt.Name)
				}
				ignoreUpgraded = true
				skipErrors = true
			} else {
				// Upgrade certain plugins
				for _, arg := range args {
					if !validation.IsSafePluginName(arg) {
						return unsafePluginNameErr(arg)
					}
					r, err := receipt.Load(paths.PluginInstallReceiptPath(arg))
					if err != nil {
						return errors.Wrapf(err, "read receipt %q", arg)
					}
					pluginNames = append(pluginNames, r.Name)
				}
			}
			client := global.SrewClient
			var nErrors int
			for _, pluginName := range pluginNames {

				res, err := client.ListPlugin(pluginName, global.PluginVersion)
				if err == nil {
					fmt.Fprintf(os.Stderr, "Upgrading plugin: %s\n", pluginName)
					err = installation.Upgrade(paths, installation.ToPlugin(res.Data))
					if ignoreUpgraded && err == installation.ErrIsAlreadyUpgraded {
						fmt.Fprintf(os.Stderr, "Skipping plugin %s, it is already on the newest version\n", pluginName)
						continue
					}
				}
				if err != nil {
					nErrors++
					if skipErrors {
						fmt.Fprintf(os.Stderr, "WARNING: failed to upgrade plugin %q, skipping (error: %v)\n", pluginName, err)
						continue
					}
					return errors.Wrapf(err, "failed to upgrade plugin %q", pluginName)
				}
				fmt.Fprintf(os.Stderr, "Upgraded plugin: %s\n", pluginName)

			}
			if nErrors > 0 {
				fmt.Fprintf(os.Stderr, "WARNING: Some plugins failed to upgrade, check logs above.\n")
			}
			return nil
		},
	}

	return upgradeCmd
}
