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

package projectmanager

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/notedownorg/notedown/pkg/providers/projects"
)

func (m *Model) submit() tea.Cmd {
	opts := make([]projects.ProjectOption, 0)
	if m.status.value != m.project.Status() {
		opts = append(opts, projects.WithStatus(m.status.value))
	}

	// Project rename requires a dedicated API call due to file renaming
	if m.text.Value() != m.project.Name() {
		return func() tea.Msg {
			p := projects.NewProjectFromProject(m.project, opts...)
			err := m.nd.RenameProject(p, m.text.Value())
			if err != nil {
				m.footer.SetMessage(fmt.Sprintf("failed to rename project: %v", err), m.ctx.Now().Add(5*time.Second), m.ctx.Theme.Red)
			}
			return nil
		}
	}

	if len(opts) == 0 {
		return nil // nothing has changed
	}

	return func() tea.Msg {
		p := projects.NewProjectFromProject(m.project, opts...)
		err := m.nd.UpdateProject(p)
		if err != nil {
			m.footer.SetMessage(fmt.Sprintf("failed to update project: %v", err), m.ctx.Now().Add(5*time.Second), m.ctx.Theme.Red)
		}
		return nil
	}
}
