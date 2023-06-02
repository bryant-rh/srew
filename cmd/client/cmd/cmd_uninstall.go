package cmd

import (
	"fmt"
	"os"

	"github.com/bryant-rh/srew/pkg/installation"
	"github.com/bryant-rh/srew/pkg/installation/validation"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

func NewCmdUninstall() *cobra.Command {
	// uninstallCmd represents the uninstall command
	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall plugins",
		Long: `Uninstall one or more plugins.

Example:
  srew uninstall NAME [NAME...]

Remarks:
  Failure to uninstall a plugin will result in an error and exit immediately.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, name := range args {
				if !validation.IsSafePluginName(name) {
					return unsafePluginNameErr(name)
				}
				klog.V(4).Infof("Going to uninstall plugin %s\n", name)
				if err := installation.Uninstall(paths, name); err != nil {
					return errors.Wrapf(err, "failed to uninstall plugin %s", name)
				}
				fmt.Fprintf(os.Stderr, "Uninstalled plugin: %s\n", name)
			}
			return nil
		},
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"remove"},
	}
	return uninstallCmd
}
