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
	"log/slog"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/notedownorg/notedown/pkg/providers/projects"
)

func (m *Model) submit() (tea.Model, tea.Cmd) {
	path := m.location.file
	name := strings.Replace(filepath.Base(path), filepath.Ext(path), "", 1)
	if err := m.nd.CreateProject(path, name, projects.Backlog); err != nil {
		slog.Error("failed to create project", "error", err)
		return m, nil
	}

	return m.ctx.Back(), nil
}
