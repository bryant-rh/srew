package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bryant-rh/srew/pkg/client"
	"github.com/bryant-rh/srew/pkg/installation"
	"github.com/bryant-rh/srew/pkg/installation/receipt"
	"github.com/bryant-rh/srew/pkg/installation/receiptsmigration"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/srew/cmd/client/global"
	"github.com/bryant-rh/srew/pkg/environment"
	"github.com/bryant-rh/srew/pkg/util"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

//定义日志级别
var (
	Verbose int
	paths   environment.Paths // srew paths used by the process
)

func NewCmd() *cobra.Command {
	// rootCmd represents the root command
	rootCmd := &cobra.Command{
		Use:   "srew",
		Short: "srew - An In-house Package Manager for macOS (or Linux)",
		Long: `srew - An In-house Package Manager for macOS (or Linux).
	You can do package management with the following command:
	"srew [command]..."`,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		PersistentPreRunE: preRun,
		Run: func(cmd *cobra.Command, args []string) {
			//_ = cmd.Help()
		},
	}
	global.SrewClient = client.NewGithubClient(global.SREW_SERVER_BASEURL)

	rootCmd.PersistentFlags().BoolVar(&global.Format, "no-format", global.Format, "If present, print output without format table")

	rootCmd.AddCommand(NewCmdSearch())
	rootCmd.AddCommand(NewCmdList())
	rootCmd.AddCommand(NewCmdInfo())
	rootCmd.AddCommand(NewCmdInstall())
	rootCmd.AddCommand(NewCmdUninstall())
	rootCmd.AddCommand(NewCmdUpgrade())
	rootCmd.AddCommand(NewCmdVersion())
	return rootCmd
}

func preRun(cmd *cobra.Command, _ []string) error {
	// check must be done before ensureDirs, to detect srew's self-installation
	if !installation.IsBinDirInPATH(paths) {
		installation.PrintWarning(os.Stderr, installation.SetupInstructions()+"\n\n")
	}

	if err := ensureDirs(paths.BasePath(),
		paths.InstallPath(),
		paths.BinPath(),
		paths.IndexBase(),
		paths.InstallReceiptsPath()); err != nil {
		klog.Fatal(err)
	}

	go func() {
		if _, disabled := os.LookupEnv("SREW_NO_UPGRADE_CHECK"); disabled ||
			global.UpgradeCheckRate < rand.Float64() { // only do the upgrade check randomly
			klog.V(1).Infof("skipping upgrade check")
			return
		}
		klog.V(1).Infof("starting upgrade check")

	}()

	// detect if receipts migration (v0.2.x->v0.3.x) is complete
	isMigrated, err := receiptsmigration.Done(paths)
	if err != nil {
		return err
	}
	if !isMigrated {
		fmt.Fprintln(os.Stderr, `This version of Srew is not supported anymore. Please manually migrate:
1. Uninstall Srew
2. Install latest Srew
3. Install the plugins you used`)
		return errors.New("srew home outdated")
	}

	if installation.IsWindows() {
		klog.V(4).Infof("detected windows, will check for old srew installations to clean up")
		err := cleanupStaleSrewInstallations()
		if err != nil {
			klog.Warningf("Failed to clean up old installations of srew (on windows).")
			klog.Warningf("You may need to clean them up manually. Error: %v", err)
		}
	}

	if len(global.SREW_SERVER_BASEURL) == 0 || len(global.SREW_SERVER_USERNAME) == 0 || len(global.SREW_SERVER_PASSWORD) == 0 {
		klog.Fatalf("请在目录:[%s] 创建配置文件: [srew.yaml], 配置SREW_SERVER_BASEURL、SREW_SERVER_USERNAME、SREW_SERVER_PASSWORD 或者配置对应环境变量\n", paths.BasePath())

	}
	if global.EnableDebug { // Enable debug mode if `--enableDebug=true` or `DEBUG=true`.
		global.SrewClient.SetDebug(true)
	}
	klog.V(4).Infof("SREW_SERVER_BASEURL: %s; SREW_SERVER_USERNAME: %s; ", global.SREW_SERVER_BASEURL, global.SREW_SERVER_USERNAME)

	authFile := fmt.Sprintf("%s/%s", paths.BasePath(), global.TokenFile)
	if !util.Exists(authFile) {
		klog.V(4).Infoln("登录生成token")

		res, err := global.SrewClient.User_Login(global.SREW_SERVER_USERNAME, global.SREW_SERVER_PASSWORD)
		if err != nil {
			klog.Fatal(err)
		}
		klog.V(4).Infof("登录生成token,并保存至文件: [%s]", authFile)

		err = ioutil.WriteFile(authFile, []byte(res.Data), 0644)
		if err != nil {
			klog.Fatal(err)
		}
		global.Token = res.Data
	} else {
		klog.V(4).Infoln("检测token是否过期")
		auth_token, err := ioutil.ReadFile(authFile)
		if err != nil {
			klog.Fatal(err)
		}
		_, err = global.SrewClient.User_VerifyToken(string(auth_token))
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				klog.V(4).Infoln("检测token已过期,重新登录生成token")
				new_res, err := global.SrewClient.User_Login(global.SREW_SERVER_USERNAME, global.SREW_SERVER_PASSWORD)
				if err != nil {
					klog.Fatal(err)
				}
				klog.V(4).Infof("生成token,并保存至文件: [%s]", authFile)
				err = ioutil.WriteFile(authFile, []byte(new_res.Data), 0644)
				if err != nil {
					klog.Fatal(err)
				}
				global.Token = new_res.Data
			} else {
				klog.Fatal(err)

			}
		} else {
			klog.V(4).Infoln("token未过期")
			global.Token = string(auth_token)

		}

	}
	//klog.V(4).Infof("token: %s", global.Token)
	klog.V(4).Infoln("执行 LoginWithToken")

	global.SrewClient.LoginWithToken(global.Token)

	return nil
}

func initLog() {
	klog.InitFlags(nil)
	rand.Seed(time.Now().UnixNano())

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	_ = flag.CommandLine.Parse([]string{}) // convince pkg/flag we parsed the flags

	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		if f.Name != "v" { // hide all glog flags except for -v
			pflag.Lookup(f.Name).Hidden = true
		}
	})
	if err := flag.Set("logtostderr", "true"); err != nil {
		fmt.Printf("can't set log to stderr %+v", err)
		os.Exit(1)
	}
}

func initConfig() {

	paths = environment.MustGetSrewPaths()
	if !util.Exists(paths.BasePath()) {
		if err := ensureDirs(paths.BasePath()); err != nil {
			klog.Fatal(err)
		}
	}

	srew_config := fmt.Sprintf("%s/srew.yaml", paths.BasePath())
	if !util.Exists(srew_config) {
		f, err := os.Create(srew_config)
		if err != nil {
			klog.Fatal(err)
		}
		f.Close()
	}
	viper.AddConfigPath(paths.BasePath())
	viper.SetConfigName("srew")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		klog.Fatalf("read config file failed! %v;\n 请在目录:[%s] 创建配置文件: [srew.yaml], 配置SREW_SERVER_BASEURL、SREW_SERVER_USERNAME、SREW_SERVER_PASSWORD 或者配置对应环境变量\n", err, paths.BasePath())

	}

	global.SREW_SERVER_BASEURL = viper.GetString("SREW_SERVER_BASEURL")
	global.SREW_SERVER_USERNAME = viper.GetString("SREW_SERVER_USERNAME")
	global.SREW_SERVER_PASSWORD = viper.GetString("SREW_SERVER_PASSWORD")

}

//init
func init() {
	initLog()
	initConfig()
}

func ensureDirs(paths ...string) error {
	for _, p := range paths {
		klog.V(4).Infof("Ensure creating dir: %q", p)

		if err := os.MkdirAll(p, 0o755); err != nil {
			return errors.Wrapf(err, "failed to ensure create directory %q", p)
		}
	}
	return nil
}

func isTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}

func cleanupStaleSrewInstallations() error {
	r, err := receipt.Load(paths.PluginInstallReceiptPath(global.SrewPluginName))
	if os.IsNotExist(err) {
		klog.V(1).Infof("could not find srew's own plugin receipt, skipping cleanup of stale srew installations")
		return nil
	} else if err != nil {
		return errors.Wrap(err, "cannot load srew's own plugin receipt")
	}
	v := r.Spec.Version

	klog.V(1).Infof("Clean up krew stale installations, current=%s", v)
	return installation.CleanupStaleSrewInstallations(paths.PluginInstallPath(global.SrewPluginName), v)
}
