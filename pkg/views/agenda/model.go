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

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/notedownorg/notedown/pkg/workspace/tasks"
	"github.com/notedownorg/task/pkg/components/statusbar"
	"github.com/notedownorg/task/pkg/components/tasklist"
	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/views/taskeditor"
)

const (
	view = "agenda"
)

func New(ctx *context.ProgramContext, t *tasks.Client) *Model {
	return &Model{
		ctx:   ctx,
		tasks: t,
		date:  time.Now(),

		keyMap:   DefaultKeyMap,
		tasklist: tasklist.New(ctx, t),
		footer:   statusbar.New(ctx, statusbar.NewMode(view, statusbar.ActionNeutral), t),
	}
}

type Model struct {
	ctx   *context.ProgramContext
	tasks *tasks.Client
	date  time.Time

	keyMap KeyMap

	tasklist *tasklist.Model
	footer   *statusbar.Model
}

func (m *Model) Init() (tea.Model, tea.Cmd) {
	_, cmd := m.ctx.Init()
	return m, cmd
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle program level key presses and events
	model, cmd := m.ctx.Update(msg)
	if model != nil {
		return model, cmd // If model is not nil, we're navigating back to the previous view
	}

	// Handle view level key presses
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.AddTask):
			return m.ctx.Navigate(m, taskeditor.NewAddModel(m.ctx, m.tasks))
		}
	}

	// Handle component events
	tl, _ := m.tasklist.Update(msg)
	m.tasklist = tl.(*tasklist.Model)

	return m, cmd
}

func (m *Model) View() string {
	horizontalPadding := 2
	verticalMargin := 1
	h := lipgloss.Height

	footer := m.footer.
		Width(m.ctx.ScreenWidth-horizontalPadding*2).
		Margin(verticalMargin, 0).
		View()

	tasklist := m.tasklist.
		Height(m.ctx.ScreenHeight-h(footer)-verticalMargin*2).
		Width(m.ctx.ScreenWidth).
		Margin(verticalMargin, 0).
		View()

	panel := lipgloss.JoinVertical(lipgloss.Top, tasklist, footer)

	return lipgloss.NewStyle().Padding(0, horizontalPadding).Render(panel)
}
