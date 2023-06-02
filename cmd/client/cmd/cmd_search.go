package cmd

import (
	"fmt"
	"runtime"

	"github.com/bryant-rh/srew/cmd/client/global"
	"github.com/bryant-rh/srew/pkg/index"
	"github.com/bryant-rh/srew/pkg/installation"
	"github.com/bryant-rh/srew/pkg/output"
	"github.com/bryant-rh/srew/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

func NewCmdSearch() *cobra.Command {
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Discover plugins",
		Long: `List plugins available on srew and search among them.
	If no arguments are provided, all plugins will be listed.
	
	Examples:
	  To list all plugins:
		srew search
	
	  To fuzzy search plugins with a keyword:
		srew search KEYWORD`,
		RunE: func(cmd *cobra.Command, args []string) error {
			klog.V(3).Infoln("Search Plugin")
			client := global.SrewClient
			if len(args) > 0 {
				global.PluginName = args[0]
			}
			receipts, err := installation.GetInstalledPluginReceipts(paths.InstallReceiptsPath())
			if err != nil {
				klog.Fatal(errors.Wrap(err, "failed to load installed plugins"))
			}

			var data [][]string
			//if global.PluginName != "" && global.AllVersion {
			if global.AllVersion {
				res, err := client.SearchPluginAllVersion(global.PluginName)
				if err != nil {
					klog.Fatal(err)
				}

				for _, v := range res.Data {
					for _, i := range v {
						platforms, err := installation.ToPlatforms(i.Platforms)
						if err != nil {
							klog.Fatal(err)
						}

						var status string
						status, _, err = statusInstalled(i.PluginName, i.Version, platforms, receipts)
						if err != nil {
							//klog.Fatal(err)
							status = fmt.Sprintf("err: %s", err)
						}
						data = append(data, []string{i.PluginName, limitString(i.ShortDescription, 60), i.Version, status})
					}
				}

				klog.V(4).Infoln("Search Plugin All Version, 输出结果")
				data = installation.SortByFirstColumn(data)
				output.Write(data, "search_all_version", global.Format, true)

			} else {

				res, err := client.SearchPlugin(global.PluginName, global.PluginVersion)
				if err != nil {
					klog.Fatal(err)
				}

				for _, v := range res.Data {
					platforms, err := installation.ToPlatforms(v.Platforms)
					if err != nil {
						klog.Fatal(err)
					}
					var status string
					var upgrade string
					status, upgrade, err = statusInstalled(v.PluginName, v.Version, platforms, receipts)
					if err != nil {
						//klog.Fatal(err)
						status = fmt.Sprintf("err: %s", err)
					}
					data = append(data, []string{v.PluginName, limitString(v.ShortDescription, 60), v.Version, status, upgrade})

				}
				klog.V(4).Infoln("Search Plugin, 输出结果")
				data = installation.SortByFirstColumn(data)
				output.Write(data, "search", global.Format, false)
			}
			return nil
		},
	}
	searchCmd.Flags().BoolVarP(&global.AllVersion, "all-version", "A", global.AllVersion, "列出所有版本")
	return searchCmd
}

func limitString(s string, length int) string {
	if len(s) > length && length > 3 {
		s = s[:length-3] + "..."
	}
	return s
}

func statusInstalled(pluginName string, pluginVersion string, platforms []index.Platform, receipts []index.Receipt) (string, string, error) {
	installed := make(map[string]bool)
	install_version := make(map[string]string)
	var status string
	var upgrade string

	for _, receipt := range receipts {
		cn := installation.DisplayName(receipt.Plugin)
		installed[cn] = true
		install_version[cn] = installation.DisplayVersion(receipt.Plugin)
	}

	if installed[pluginName] {
		if global.AllVersion {
			if install_version[pluginName] == pluginVersion {
				status = util.GreenColor("yes")

			}else{
				status = util.YellowColor("no")

			}

		} else {
			status = util.GreenColor("yes")

		}

	} else if _, ok, err := installation.GetMatchingPlatform(platforms); err != nil {
		return "", "", errors.Wrapf(err, "failed to get the matching platform for plugin %s", pluginName)
	} else if ok {
		status = util.YellowColor("no")
	} else {
		status = util.RedColor(fmt.Sprintf("unavailable on %v/%v", runtime.GOOS, runtime.GOARCH))
	}

	upgrade = statusUpgrade(install_version[pluginName], pluginVersion)

	return status, upgrade, nil
}

func statusUpgrade(current_version string, latest_version string) string {
	var status string
	res := util.CompareStrVer(current_version, latest_version)

	switch res {
	case 0:
		status = util.GreenColor("no")
	case 1:
		status = util.GreenColor("no")
	case 2:
		status = util.RedColor("yes")
	default:
		status = ""
	}
	return status
}
