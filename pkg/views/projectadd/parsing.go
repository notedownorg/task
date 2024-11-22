// Copyright 2024 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package projectadd

import (
	"path/filepath"
	"strings"
)

func (m *Model) parse() {
	txt := m.text.Value()

	if isRelativePath(txt) {
		m.location.SetLocation(filepath.Clean(txt))
	} else {
		m.location.SetLocation(m.nd.NewProjectLocation(txt))
	}

	m.text.IsValid = isValidPath(m.location.file)

}

func isValidPath(path string) bool {
	if path == "" {
		return false
	}

	return filepath.Ext(path) == ".md"
}

func isRelativePath(path string) bool {
	if filepath.IsAbs(path) {
		return false
	}

	if filepath.Ext(path) == ".md" {
		return true
	}

	if strings.ContainsAny(path, "./\\") {
		return true
	}

	return false
}
