package cmd

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/bryant-rh/srew/cmd/client/global"
	"github.com/bryant-rh/srew/internal/model"
	"github.com/bryant-rh/srew/pkg/installation"
	"github.com/bryant-rh/srew/pkg/util"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

func NewCmdInfo() *cobra.Command {
	infoCmd := &cobra.Command{
		Use:     "info",
		Short:   "Show information about an available plugin",
		Long:    `Show detailed information about an available plugin.`,
		Example: `  srew info PLUGIN `,
		RunE: func(cmd *cobra.Command, args []string) error {
			klog.V(3).Infoln("Info Plugin")
			client := global.SrewClient
			if len(args) > 0 {
				global.PluginName = args[0]
			}
			res, err := client.ListPlugin(global.PluginName, global.PluginVersion)
			if err != nil {
				klog.Fatal(err)
			}
			for _, v := range res.Data {
				printPluginInfo(os.Stdout, v)
			}
			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	infoCmd.Flags().StringVar(&global.PluginVersion, "version", global.PluginVersion, "可指定插件版本进行查看")
	return infoCmd
}

func printPluginInfo(out io.Writer, plugin model.Detail) {
	fmt.Fprintf(out, "%s\t: %s\n", util.GreenColor("NAME"), util.CyanColor(plugin.PluginName))
	if plugin.Version != "" {
		fmt.Fprintf(out, "%s\t: %s\n", util.GreenColor("VERSION"), util.CyanColor(plugin.Version))
	}
	platforms, err := installation.ToPlatforms(plugin.Platforms)
	if err != nil {
		klog.Fatal(err)
	}
	if platform, ok, err := installation.GetMatchingPlatform(platforms); err == nil && ok {
		if platform.URI != "" {
			fmt.Fprintf(out, "%s\t: %s\n", util.GreenColor("URI"), util.CyanColor(platform.URI))
			fmt.Fprintf(out, "%s\t: %s\n", util.GreenColor("SHA256"), util.CyanColor(platform.Sha256))
		}
	}

	if plugin.Homepage != "" {
		fmt.Fprintf(out, "%s: %s\n", util.GreenColor("HOMEPAGE"), util.CyanColor(plugin.Homepage))
	}
	if plugin.Description != "" {
		fmt.Fprintf(out, "%s: \n%s\n", util.GreenColor("DESCRIPTION"), util.CyanColor(plugin.Description))
	}
	if plugin.Caveats != "" {
		fmt.Fprintf(out, "%s:\n%s\n", util.GreenColor("CAVEATS"), util.CyanColor(indent(plugin.Caveats)))
	}
}

// indent converts strings to an indented format ready for printing.
// Example:
//
//     \
//      | This plugin is great, use it with great care.
//      | Also, plugin will require the following programs to run:
//      |  * jq
//      |  * base64
//     /
func indent(s string) string {
	out := "\\\n"
	s = strings.TrimRightFunc(s, unicode.IsSpace)
	out += regexp.MustCompile("(?m)^").ReplaceAllString(s, " | ")
	out += "\n/"
	return out
}
