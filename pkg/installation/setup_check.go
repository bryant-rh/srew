// Copyright 2020 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package installation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/klog/v2"

	"github.com/bryant-rh/srew/pkg/environment"
)

const (
	instructionWindows = `To be able to run plugins, you need to add the
"%%USERPROFILE%%\.srew\bin" directory to your PATH environment variable
and restart your shell.`
	instructionUnixTemplate = `To be able to run srew plugins, you need to add
the following to your %s

and restart your shell.`
	instructionZsh = `~/.zshrc:

    export PATH="${SREW_ROOT:-$HOME/.srew}/bin:$PATH"`
	instructionBash = `~/.bash_profile or ~/.bashrc:

    export PATH="${SREW_ROOT:-$HOME/.srew}/bin:$PATH"`
	instructionFish = `config.fish:

    set -q SREW_ROOT; and set -gx PATH $PATH $SREW_ROOT/.srew/bin; or set -gx PATH $PATH $HOME/.srew/bin`
	instructionGeneric = `~/.bash_profile, ~/.bashrc, or ~/.zshrc:

    export PATH="${SREW_ROOT:-$HOME/.srew}/bin:$PATH"`
)

func IsBinDirInPATH(paths environment.Paths) bool {
	// For the first run we don't want to show a warning.
	_, err := os.Stat(paths.BasePath())
	if err != nil {
		klog.V(4).Info("Assuming this is the first run")
		return os.IsNotExist(err)
	}

	binPath := paths.BinPath()
	for _, dirInPATH := range filepath.SplitList(os.Getenv("PATH")) {
		normalizedDirInPATH, err := filepath.Abs(dirInPATH)
		if err != nil {
			klog.Warningf("Cannot get absolute path: %v, %v", normalizedDirInPATH, err)
			continue
		}
		if normalizedDirInPATH == binPath {
			return true
		}
	}
	return false
}

func SetupInstructions() string {
	// if installation.IsWindows() {
	// 	return instructionWindows
	// }

	var instruction string
	switch shell := os.Getenv("SHELL"); {
	case strings.HasSuffix(shell, "/zsh"):
		instruction = instructionZsh
	case strings.HasSuffix(shell, "/bash"):
		instruction = instructionBash
	case strings.HasSuffix(shell, "/fish"):
		instruction = instructionFish
	default:
		instruction = instructionGeneric
	}
	return fmt.Sprintf(instructionUnixTemplate, instruction)
}
