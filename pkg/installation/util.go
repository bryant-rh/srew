// Copyright 2019 The Kubernetes Authors.
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
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"sort"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/srew/cmd/client/global"
	"github.com/bryant-rh/srew/internal/model"
	"github.com/bryant-rh/srew/pkg/index"
	"github.com/bryant-rh/srew/pkg/installation/receipt"
)

// InstalledPluginsFromIndex returns a list of all install plugins from a particular index.
func InstalledPluginsFromIndex(receiptsDir, indexName string) ([]index.Receipt, error) {
	var out []index.Receipt
	receipts, err := GetInstalledPluginReceipts(receiptsDir)
	if err != nil {
		return nil, err
	}
	out = append(out, receipts...)
	return out, nil
}

// GetInstalledPluginReceipts returns a list of receipts.
func GetInstalledPluginReceipts(receiptsDir string) ([]index.Receipt, error) {
	files, err := filepath.Glob(filepath.Join(receiptsDir, "*"+global.ManifestExtension))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to glob receipts directory (%s) for manifests", receiptsDir)
	}
	out := make([]index.Receipt, 0, len(files))
	for _, f := range files {
		r, err := receipt.Load(f)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse plugin install receipt %s", f)
		}
		out = append(out, r)
		klog.V(4).Infof("parsed receipt for %s: version=%s", r.GetObjectMeta().GetName(), r.Spec.Version)

	}
	return out, nil
}

func DisplayName(p index.Plugin) string {
	return p.Name
}

func DisplayVersion(p index.Plugin) string {
	return p.Spec.Version
}

func SortByFirstColumn(rows [][]string) [][]string {
	sort.Slice(rows, func(a, b int) bool {
		return rows[a][0] < rows[b][0]
	})
	return rows
}

func ToPlatforms(s []string) ([]index.Platform, error) {
	var platforms []index.Platform
	for _, p := range s {
		platform := index.Platform{}
		err := json.Unmarshal([]byte(p), &platform)
		if err != nil {
			return nil, err
		}

		platforms = append(platforms, platform)
	}
	return platforms, nil
}

func ToPlugin(m []model.Detail) index.Plugin {

	var plugin = index.Plugin{}
	for _, p := range m {
		plugin.APIVersion = global.CurrentAPIVersion
		plugin.Kind = global.PluginKind
		plugin.Name = p.PluginName
		plugin.Spec.Version = p.Version
		plugin.Spec.Homepage = p.Homepage
		plugin.Spec.ShortDescription = p.ShortDescription
		plugin.Spec.Description = p.Description
		plugin.Spec.Caveats = p.Caveats
		platforms, err := ToPlatforms(p.Platforms)
		if err != nil {
			klog.Fatal(err)
		}
		plugin.Spec.Platforms = platforms

	}
	return plugin
}

func ToPluginDetail(m map[string]interface{}) error {
	details := map[string][]model.Detail{}
	// convert map to json
	jsonString, _ := json.Marshal(m)
	//fmt.Println(string(jsonString))

	// convert json to struct
	json.Unmarshal(jsonString, &details)
	fmt.Println(details)

	return nil
}

func PrintWarning(w io.Writer, format string, a ...interface{}) {
	color.New(color.FgRed, color.Bold).Fprint(w, "WARNING: ")
	fmt.Fprintf(w, format, a...)
}
