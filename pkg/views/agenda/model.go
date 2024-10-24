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

package agenda

import (
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/notedownorg/notedown/pkg/workspace/tasks"
	"github.com/notedownorg/task/pkg/components/statusbar"
	"github.com/notedownorg/task/pkg/context"
)

func New(ctx *context.ProgramContext, t *tasks.Client) *Model {
	return &Model{
		ctx:   ctx,
		tasks: t,
		date:  time.Now(),
	}
}

type Model struct {
	ctx   *context.ProgramContext
	tasks *tasks.Client
	date  time.Time
}

func (m *Model) Init() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	ctx, cmd := m.ctx.Update(msg)
	m.ctx = ctx.(*context.ProgramContext)
	return m, cmd
}

func (m *Model) View() string {
	stats := statusbar.Stats{
		Tasks:     len(m.tasks.ListTasks(tasks.FetchAllTasks())),
		Projects:  len(m.tasks.ListDocuments(tasks.FetchAllDocuments(), tasks.FilterByDocumentType("project"))),
		Documents: len(m.tasks.ListDocuments(tasks.FetchAllDocuments())),
	}
	return statusbar.Render(m.ctx, "agenda", stats)
}
