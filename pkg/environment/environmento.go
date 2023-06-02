package environment

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/bryant-rh/srew/cmd/client/global"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

type Paths struct {
	base string
	tmp  string
}

// HomeDir returns the home directory for the current user
func HomeDir() string {
	if runtime.GOOS == "windows" {

		// First prefer the HOME environmental variable
		if home := os.Getenv("HOME"); len(home) > 0 {
			if _, err := os.Stat(home); err == nil {
				return home
			}
		}
		if homeDrive, homePath := os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"); len(homeDrive) > 0 && len(homePath) > 0 {
			homeDir := homeDrive + homePath
			if _, err := os.Stat(homeDir); err == nil {
				return homeDir
			}
		}
		if userProfile := os.Getenv("USERPROFILE"); len(userProfile) > 0 {
			if _, err := os.Stat(userProfile); err == nil {
				return userProfile
			}
		}
	}
	return os.Getenv("HOME")
}

//MustGetSrewPaths
func MustGetSrewPaths() Paths {
	base := filepath.Join(HomeDir(), ".srew")
	if fromEnv := os.Getenv("SREW_ROOT"); fromEnv != "" {
		base = fromEnv
		klog.V(4).Infof("using environment override SREW_ROOT=%s", fromEnv)
	}
	base, err := filepath.Abs(base)
	if err != nil {
		panic(errors.Wrap(err, "cannot get absolute path"))
	}
	return NewPaths(base)
}

//NewPaths
func NewPaths(base string) Paths {
	return Paths{base: base, tmp: os.TempDir()}
}

// BasePath returns srew base directory.
func (p Paths) BasePath() string { return p.base }

// e.g. {BasePath}/bin
func (p Paths) BinPath() string { return filepath.Join(p.base, "bin") }

// e.g. {BasePath}/store
func (p Paths) InstallPath() string { return filepath.Join(p.base, "store") }

// index and custom ones.
func (p Paths) IndexBase() string {
	return filepath.Join(p.base, "index")
}

// IndexPath returns the directory where a plugin index repository is cloned.
// e.g. {BasePath}/index/default or {BasePath}/index
func (p Paths) IndexPath(name string) string {
	return filepath.Join(p.base, "index", name)
}

// InstallReceiptsPath returns the base directory where plugin receipts are stored.
//
// e.g. {BasePath}/receipts
func (p Paths) InstallReceiptsPath() string { return filepath.Join(p.base, "receipts") }

// e.g. {InstallPath}/{version}/{..files..}
func (p Paths) PluginInstallPath(plugin string) string {
	return filepath.Join(p.InstallPath(), plugin)
}

// PluginInstallReceiptPath returns the path to the install receipt for plugin.
//
// e.g. {InstallReceiptsPath}/{plugin}.yaml
func (p Paths) PluginInstallReceiptPath(plugin string) string {
	return filepath.Join(p.InstallReceiptsPath(), plugin+global.ManifestExtension)
}

// PluginVersionInstallPath returns the path to the specified version of specified
// plugin.
//
// e.g. {PluginInstallPath}/{plugin}/{version}
func (p Paths) PluginVersionInstallPath(plugin, version string) string {
	return filepath.Join(p.InstallPath(), plugin, version)
}

// Realpath evaluates symbolic links. If the path is not a symbolic link, it
// returns the cleaned path. Symbolic links with relative paths return error.
func Realpath(path string) (string, error) {
	s, err := os.Lstat(path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to stat the currently executed path (%q)", path)
	}

	if s.Mode()&os.ModeSymlink != 0 {
		if path, err = os.Readlink(path); err != nil {
			return "", errors.Wrap(err, "failed to resolve the symlink of the currently executed version")
		}
		if !filepath.IsAbs(path) {
			return "", errors.Errorf("symbolic link is relative (%s)", path)
		}
	}
	return filepath.Clean(path), nil
}

//SetCofnigFile
func SetCofnigFile(path string) error {
	viper.SetConfigFile(path)
	return nil
}
