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

package taskadd

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/notedownorg/notedown/pkg/workspace/tasks"
	"github.com/notedownorg/task/pkg/components/statusbar"
	"github.com/notedownorg/task/pkg/context"
)

type Model struct {
	ctx   *context.ProgramContext
	tasks *tasks.Client

	footer *statusbar.Model
}

func NewModel(ctx *context.ProgramContext, t *tasks.Client) *Model {
	return &Model{
		ctx:   ctx,
		tasks: t,

		footer: statusbar.New(ctx, statusbar.NewMode("add task", statusbar.ActionCreate), t),
	}
}

func (m *Model) Init() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle program level key presses and events
	model, cmd := m.ctx.Update(msg)

	// If model is not nil, we're navigating back to the previous view
	if model != nil {
		return model, cmd
	}

	// Handle view level key presses

	// Handle component events

	return m, cmd
}

func (m *Model) View() string {
	horizontalPadding := 2
	verticalMargin := 1

	footer := m.footer.
		Width(m.ctx.ScreenWidth-horizontalPadding*2).
		Margin(verticalMargin, 0).
		View()

	panel := lipgloss.JoinVertical(lipgloss.Top, footer)

	return lipgloss.NewStyle().Padding(0, horizontalPadding).Render(panel)
}
