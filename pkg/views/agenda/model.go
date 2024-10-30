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
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/notedownorg/notedown/pkg/ast"
	"github.com/notedownorg/notedown/pkg/workspace/tasks"
	"github.com/notedownorg/task/pkg/components/groupedlist"
	"github.com/notedownorg/task/pkg/components/statusbar"
	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/views/taskeditor"
)

const (
	view = "agenda"
)

func New(ctx *context.ProgramContext, t *tasks.Client) *Model {
	m := &Model{
		ctx:   ctx,
		tasks: t,

		keyMap: DefaultKeyMap,
		date:   time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC),

		tasklist: groupedlist.New[ast.Task](groupedlist.WithRenderers(rendererFuncs(ctx.Theme))),
		footer:   statusbar.New(ctx, statusbar.NewMode(view, statusbar.ActionNeutral), t),
	}
	m.tasklist.SetGroups(m.visibleTasks())
	return m
}

type Model struct {
	ctx   *context.ProgramContext
	tasks *tasks.Client

	keyMap KeyMap
	date   time.Time

	tasklist *groupedlist.Model[ast.Task]
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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.AddTask):
			return m.ctx.Navigate(m, taskeditor.NewAddModel(m.ctx, m.tasks))
		case key.Matches(msg, m.keyMap.NextDay):
			m.updateDate(m.date.AddDate(0, 0, 1))
		case key.Matches(msg, m.keyMap.PrevDay):
			m.updateDate(m.date.AddDate(0, 0, -1))
		case key.Matches(msg, m.keyMap.ResetDate):
			m.date = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
		case key.Matches(msg, m.keyMap.CursorUp):
			m.tasklist.MoveUp(1)
		case key.Matches(msg, m.keyMap.CursorDown):
			m.tasklist.MoveDown(1)
		}
	}

	return m, cmd
}

func (m *Model) updateDate(date time.Time) {
	m.date = date
	m.tasklist.SetGroups(m.visibleTasks())
}

func (m *Model) View() string {
	gap := ""
	horizontalPadding := 2
	verticalPadding := 1
	h := lipgloss.Height

	header := fmt.Sprintf("← %v →", humanizeDate(m.date, time.Now()))

	footer := m.footer.
		Width(m.ctx.ScreenWidth - horizontalPadding*2).
		View()

	tasklist := m.tasklist.
		Height(m.ctx.ScreenHeight - h(footer) - h(header) - verticalPadding*2 - 2). // -2 for the gaps
		Width(m.ctx.ScreenWidth - horizontalPadding*2).
		View()

	panel := lipgloss.JoinVertical(lipgloss.Top, header, gap, tasklist, gap, footer)

	return lipgloss.NewStyle().Padding(verticalPadding, horizontalPadding).Render(panel)
}

func humanizeDate(date time.Time, relativeTo time.Time) string {
	rel := time.Date(relativeTo.Year(), relativeTo.Month(), relativeTo.Day(), 0, 0, 0, 0, time.UTC)
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	diff := date.Sub(rel) / (24 * time.Hour)

	switch diff {
	case -6, -5, -4, -3, -2:
		return fmt.Sprintf("Last %s", date.Weekday().String())
	case -1:
		return "Yesterday"
	case 0:
		return "Today"
	case 1:
		return "Tomorrow"
	case 2, 3, 4, 5, 6:
		return date.Weekday().String()
	}

	if date.Year() == relativeTo.Year() {
		return date.Format("Monday, January 2")
	}
	return date.Format("Monday, January 2, 2006")
}
