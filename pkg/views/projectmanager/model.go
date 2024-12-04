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
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/notedownorg/notedown/pkg/providers/projects"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/notedownorg/task/pkg/components/groupedlist"
	"github.com/notedownorg/task/pkg/components/statusbar"
	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/listeners"
	"github.com/notedownorg/task/pkg/notedown"
	"github.com/notedownorg/task/pkg/styling/tasklists"
)

type Model struct {
	ctx *context.ProgramContext
	nd  notedown.Client

	keyMap  KeyMap
	project projects.Project

	text      *Text
	status    *Status
	tasklist  *groupedlist.Model[tasks.Task]
	completed *groupedlist.Model[tasks.Task]
	footer    *statusbar.Model
}

func New(ctx *context.ProgramContext, nd notedown.Client, project projects.Project) *Model {
	m := &Model{
		ctx:       ctx,
		nd:        nd,
		keyMap:    DefaultKeyMap,
		project:   project,
		status:    NewStatus(ctx, project.Status()),
		text:      NewText(ctx, project.Name()),
		completed: groupedlist.New(groupedlist.WithRenderers(tasklists.CompletedRenderers(ctx.Theme))),
		footer:    statusbar.New(ctx, statusbar.NewMode("manage project", statusbar.ActionNeutral), nd),
	}
	m.tasklist = groupedlist.New(groupedlist.WithRenderers(tasklists.MainRenderers(ctx.Theme, ctx.Now))).Focus()
	m.updateTasks()

	return m
}

func (m *Model) Init() (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if model, command := m.text.Init(); command != nil {
		m.text, cmd = model.(*Text), tea.Batch(cmd, command)
	}
	if model, command := m.status.Init(); command != nil {
		m.status, cmd = model.(*Status), tea.Batch(cmd, command)
	}

	return m, cmd
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.ToggleFocus):
			cmd = tea.Batch(m.toggleFocus())
		case key.Matches(msg, m.keyMap.CursorUp):
			m.moveUp(1)
		case key.Matches(msg, m.keyMap.CursorDown):
			m.moveDown(1)
		}
	}

	// Handle component events
	m.status.Update(msg)
	m.text.Update(msg)

	// If the task client has emitted an event, refresh the tasks
	if _, ok := msg.(listeners.TaskEvent); ok {
		m.updateTasks()
	}

	// If we're being navigated back to, refresh the tasks
	// This is mostly just in case we miss the task event on an add/edit/etc
	if _, ok := msg.(context.NavigationEvent); ok {
		m.updateTasks()
	}

	// Handle program level key presses and events
	model, command := m.ctx.Update(msg)
	if model != nil {
		// If a model is not nil, we're navigating to a new view, submit any changes
		return model, tea.Batch(command, cmd, m.submit())
	}
	cmd = tea.Batch(cmd, command)
	return m, cmd
}

// Status -> Text -> TaskList -> Completed
func (m *Model) toggleFocus() tea.Cmd {
	if m.status.focused {
		m.status.Blur()
		m.text.Focus()
		return m.submit() // write any changes to the status
	}
	if m.text.focused {
		m.text.Blur()
		m.tasklist.Focus()
		return m.submit() // write any changes to the name
	}
	if m.tasklist.Focused() {
		m.tasklist.Blur()
		m.completed.Focus()
		return nil
	}
	if m.completed.Focused() {
		m.completed.Blur()
		m.status.Focus()
		return nil
	}
	return nil
}

func (m *Model) moveUp(n int) {
	if m.tasklist.Focused() {
		m.tasklist.MoveUp(n)
	}
	if m.completed.Focused() {
		m.completed.MoveUp(n)
	}
}

func (m *Model) moveDown(n int) {
	if m.tasklist.Focused() {
		m.tasklist.MoveDown(n)
	}
	if m.completed.Focused() {
		m.completed.MoveDown(n)
	}
}
