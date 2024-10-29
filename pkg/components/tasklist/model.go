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

package tasklist

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/notedownorg/notedown/pkg/ast"
	"github.com/notedownorg/notedown/pkg/workspace/tasks"
	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/model"
)

type Model struct {
	base   model.Base
	ctx    *context.ProgramContext
	keyMap KeyMap

	tasks *tasks.Client

	// this may need to be an array of linked list ints(?) in the future to
	// suppport selecting multiple tasks and subtasks
	selected int
}

func New(ctx *context.ProgramContext, t *tasks.Client) *Model {
	return &Model{
		ctx:    ctx,
		tasks:  t,
		keyMap: DefaultKeyMap,
	}
}

func (m *Model) Init() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Up):
			if m.selected == 0 {
				return m, nil
			}
			m.selected--
		case key.Matches(msg, m.keyMap.Down):
			m.selected++ // we don't know if this is oob until render time
		}
	}
	return m, nil
}

func (m *Model) Height(height int) *Model {
	m.base.Height(height)
	return m
}

func (m *Model) Width(width int) *Model {
	m.base.Width(width)
	return m
}

func (m *Model) Margin(margin ...int) *Model {
	m.base.Margin(margin...)
	return m
}

func (m Model) listStyle() lipgloss.Style {
	return m.base.NewStyle()
}

func (m *Model) visibleTasks() []ast.Task {
	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.Add(-1)
	tomorrow := today.AddDate(0, 0, 1).Add(-1)

	overdue := m.tasks.ListTasks(
		tasks.FetchAllTasks(),
		tasks.WithFilters(
			tasks.FilterByStatus(ast.Todo, ast.Doing, ast.Blocked),
			tasks.FilterByDueDate(nil, &tomorrow),
		),
		tasks.WithSorters(
			tasks.SortByStatus(tasks.AgendaOrder()),
			tasks.SortByPriority(),
		),
	)

	completed := m.tasks.ListTasks(
		tasks.FetchAllTasks(),
		tasks.WithFilters(
			tasks.FilterByStatus(ast.Done),
			tasks.FilterByCompletedDate(&yesterday, &tomorrow),
		),
		tasks.WithSorters(tasks.SortByAlphabetical()),
	)
	return append(overdue, completed...)
}

func (m *Model) View() string {
	t := m.visibleTasks()
	if len(t) == 0 {
		return m.listStyle().Render("No tasks")
	}

	var b strings.Builder
	for m.selected >= len(t) {
		m.selected--
	}

	currentGroup := t[0].Status()
	for i, task := range t {
		if task.Status() != currentGroup {
			currentGroup = task.Status()
			b.WriteString("\n")
		}
		b.WriteString(renderTask(m.ctx.Theme, task, i == m.selected))
		b.WriteString("\n")
	}

	return m.listStyle().
		Render(b.String())
}
