package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/bryant-rh/srew/cmd/client/global"
	"github.com/bryant-rh/srew/pkg/index"
	"github.com/bryant-rh/srew/pkg/installation"
	"github.com/bryant-rh/srew/pkg/installation/scanner"
	"github.com/bryant-rh/srew/pkg/installation/validation"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

func NewCmdInstall() *cobra.Command {

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install plugins",
		Long: `Install one or multiple plugins.

Examples:
  To install one or multiple plugins, run:
    srew install NAME [NAME...]

  To install plugins from a newline-delimited file, run:
    srew install < file.txt

Remarks:
  If a plugin is already installed, it will be skipped.
  Failure to install a plugin will not stop the installation of other plugins.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var pluginNames = make([]string, len(args))
			copy(pluginNames, args)

			if !isTerminal(os.Stdin) && (len(pluginNames) != 0 || global.Manifest != "") {
				fmt.Fprintln(os.Stderr, "WARNING: Detected stdin, but discarding it because of --manifest or args")
			}

			if !isTerminal(os.Stdin) && (len(pluginNames) == 0 && global.Manifest == "") {
				fmt.Fprintln(os.Stderr, "Reading plugin names via stdin")
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Split(bufio.ScanLines)
				for scanner.Scan() {
					if name := scanner.Text(); name != "" {
						pluginNames = append(pluginNames, name)
					}
				}
			}

			if global.Manifest != "" && global.ManifestURL != "" {
				return errors.New("cannot specify --manifest and --manifest-url at the same time")
			}

			if len(pluginNames) != 0 && (global.Manifest != "" || global.ManifestURL != "") {
				return errors.New("must specify either specify either plugin names (via positional arguments or STDIN), or --manifest/--manifest-url; not both")
			}

			client := global.SrewClient

			var install []index.Plugin
			for _, pluginName := range pluginNames {
				//		indexName, pluginName := pathutil.CanonicalPluginName(name)

				if !validation.IsSafePluginName(pluginName) {
					return unsafePluginNameErr(pluginName)
				}

				res, err := client.ListPlugin(pluginName, global.PluginVersion)
				if err != nil {
					klog.Fatal(err)
				}

				install = append(install, installation.ToPlugin(res.Data))

			}

			if global.Manifest != "" {
				plugin, err := scanner.ReadPluginFromFile(global.Manifest)
				if err != nil {
					return errors.Wrap(err, "failed to load plugin manifest from file")
				}
				install = append(install, plugin)

			} else if global.ManifestURL != "" {
				plugin, err := readPluginFromURL(global.ManifestURL)
				if err != nil {
					return errors.Wrap(err, "failed to read plugin manifest file from url")
				}
				install = append(install, plugin)
			}

			if len(install) == 0 {
				return cmd.Help()
			}

			for _, p := range install {
				klog.V(2).Infof("Will install plugin: %s\n", p.Name)
			}

			var failed []string
			var returnErr error
			for _, plugin := range install {
				fmt.Fprintf(os.Stderr, "Installing plugin: %s\n", plugin.Name)
				err := installation.Install(paths, plugin, installation.InstallOpts{
					ArchiveFileOverride: global.ArchiveFileOverride,
				})

				if err == installation.ErrIsAlreadyInstalled {
					klog.Warningf("Skipping plugin %q, it is already installed", plugin.Name)
					continue
				}
				if err != nil {
					klog.Warningf("failed to install plugin %q: %v", plugin.Name, err)
					if returnErr == nil {
						returnErr = err
					}
					failed = append(failed, plugin.Name)
					continue
				}
				fmt.Fprintf(os.Stderr, "Installed plugin: %s\n", plugin.Name)
				output := fmt.Sprintf("Use this plugin:\n\t %s\n", plugin.Name)
				if plugin.Spec.Homepage != "" {
					output += fmt.Sprintf("Documentation:\n\t%s\n", plugin.Spec.Homepage)
				}
				if plugin.Spec.Caveats != "" {
					output += fmt.Sprintf("Caveats:\n%s\n", indent(plugin.Spec.Caveats))
				}
				fmt.Fprintln(os.Stderr, indent(output))

			}
			if len(failed) > 0 {
				return errors.Wrapf(returnErr, "failed to install some plugins: %+v", failed)
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if global.Manifest != "" {
				klog.V(4).Infof("--manifest specified, not ensuring plugin ")
				return nil
			}
			return nil
		},
	}

	installCmd.Flags().StringVar(&global.PluginVersion, "version", global.PluginVersion, "可指定插件版本进行安装")
	installCmd.Flags().StringVar(&global.Manifest, "manifest", "", "(Development-only) specify local plugin manifest file")
	installCmd.Flags().StringVar(&global.ManifestURL, "manifest-url", "", "(Development-only) specify plugin manifest file from url")

	return installCmd

}

func unsafePluginNameErr(n string) error { return errors.Errorf("plugin name %q not allowed", n) }

func readPluginFromURL(url string) (index.Plugin, error) {
	klog.V(4).Infof("downloading manifest from url %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return index.Plugin{}, errors.Wrapf(err, "request to url failed (%s)", url)
	}
	klog.V(4).Infof("manifest downloaded from url, status=%v headers=%v", resp.Status, resp.Header)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return index.Plugin{}, errors.Errorf("unexpected status code (http %d) from url", resp.StatusCode)
	}
	return scanner.ReadPlugin(resp.Body)
}
