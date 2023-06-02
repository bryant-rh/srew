package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/bryant-rh/srew/cmd/client/global"
	"github.com/bryant-rh/srew/pkg/installation"
	"github.com/bryant-rh/srew/pkg/output"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCmdList() *cobra.Command {

	// listCmd represents the list command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List installed plugins",
		Long: `Show a list of installed  plugins and their versions.

Remarks:
  Redirecting the output of this command to a program or file will only print
  the names of the plugins installed. This output can be piped back to the
  "install" command.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			receipts, err := installation.GetInstalledPluginReceipts(paths.InstallReceiptsPath())
			if err != nil {
				return errors.Wrap(err, "failed to find all installed versions")
			}

			// return sorted list of plugin names when piped to other commands or file
			if !isTerminal(os.Stdout) {
				var names []string
				for _, r := range receipts {
					names = append(names, installation.DisplayName(r.Plugin))
				}
				sort.Strings(names)
				fmt.Fprintln(os.Stdout, strings.Join(names, "\n"))
				return nil
			}

			// print table
			var data [][]string
			for _, r := range receipts {
				data = append(data, []string{installation.DisplayName(r.Plugin), r.Spec.Version})
			}
			data = installation.SortByFirstColumn(data)
			output.Write(data, "list", global.Format, false)
			return nil
		},
	}

	return listCmd
}
