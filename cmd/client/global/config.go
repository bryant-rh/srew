package global

import "github.com/bryant-rh/srew/pkg/client"

//全局配置
var (
	SREW_SERVER_BASEURL  string
	SREW_SERVER_USERNAME string
	SREW_SERVER_PASSWORD string
	Token                string
	//latestTag            = ""
	TokenFile  = "token.json"
	SrewClient *client.SrewClient
)

const (
	CurrentAPIVersion = "srew.sensors.com/v1alpha2"
	PluginKind        = "Plugin"
	ManifestExtension = ".yaml"
	SrewPluginName    = "srew" // plugin name of srew itself
)

//命令参数
var (
	EnableDebug           bool
	Format                bool
	AllVersion            bool
	PluginName            string
	PluginVersion         string
	ArchiveFileOverride   string
	Manifest, ManifestURL string
)

const (
	UpgradeNotification = "A newer version of srew is available (%s -> %s).\nRun \"srew upgrade\" to get the newest version!\n"

	// upgradeCheckRate is the percentage of srew runs for which the upgrade check is performed.
	UpgradeCheckRate = 0.1
)
