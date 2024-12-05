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
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/notedownorg/task/pkg/components/groupedlist"
	"github.com/notedownorg/task/pkg/components/statusbar"
	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/listeners"
	"github.com/notedownorg/task/pkg/notedown"
	"github.com/notedownorg/task/pkg/styling/tasklists"
	"github.com/notedownorg/task/pkg/views/taskeditor"
	"github.com/notedownorg/task/pkg/views/taskreschedule"
)

const (
	view = "agenda"
)

func HandleNew(nd notedown.Client) context.GlobalKeyHandler {
	return func(ctx *context.ProgramContext, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
		key := msg.Key()
		if key.Mod == tea.ModCtrl && key.Code == 'a' {
			return ctx.Navigate(New(ctx, nd))
		}
		return nil, nil
	}
}

func New(ctx *context.ProgramContext, nd notedown.Client) *Model {
	date := time.Date(ctx.Now().Year(), ctx.Now().Month(), ctx.Now().Day(), 0, 0, 0, 0, ctx.Now().Location())
	m := &Model{
		ctx: ctx,
		nd:  nd,

		keyMap: DefaultKeyMap,
		date:   date,

		completed: groupedlist.New(groupedlist.WithRenderers(tasklists.CompletedRenderers(ctx.Theme))),
		footer:    statusbar.New(ctx, statusbar.NewMode(view, statusbar.ActionNeutral), nd),
	}
	m.tasklist = groupedlist.New(groupedlist.WithRenderers(tasklists.MainRenderers(ctx.Theme, func() time.Time { return m.date }))).Focus()
	m.updateTasks()
	return m
}

type Model struct {
	ctx *context.ProgramContext
	nd  notedown.Client

	keyMap KeyMap
	date   time.Time

	tasklist  *groupedlist.Model[tasks.Task]
	completed *groupedlist.Model[tasks.Task]
	footer    *statusbar.Model
}

func (m *Model) Init() (tea.Model, tea.Cmd) {
	_, cmd := m.ctx.Init()
	return m, cmd
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		// Internal to the agenda view
		case key.Matches(msg, m.keyMap.TogglePanels):
			m.togglePanels()
		case key.Matches(msg, m.keyMap.NextDay):
			m.updateDate(m.date.AddDate(0, 0, 1))
		case key.Matches(msg, m.keyMap.PrevDay):
			m.updateDate(m.date.AddDate(0, 0, -1))
		case key.Matches(msg, m.keyMap.ResetDate):
			m.date = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
		case key.Matches(msg, m.keyMap.CursorUp):
			m.moveUp(1)
		case key.Matches(msg, m.keyMap.CursorDown):
			m.moveDown(1)

		// Navigation
		case key.Matches(msg, m.keyMap.AddTask):
			return m.ctx.Navigate(taskeditor.New(
				m.ctx,
				m.nd,
				taskeditor.WithAddToDaily(tasks.Todo, fmt.Sprintf(" due:%s", m.date.Format("2006-01-02")), m.date),
			))
		case key.Matches(msg, m.keyMap.EditTask):
			if selected := m.selectedTask(); selected != nil {
				return m.ctx.Navigate(taskeditor.New(
					m.ctx,
					m.nd,
					taskeditor.WithEdit(*selected, m.date),
				))
			}

		// Other Task operations
		case key.Matches(msg, m.keyMap.RescheduleTask):
			if selected := m.selectedTask(); selected != nil {
				return m.ctx.Navigate(taskreschedule.New(m.ctx, m.nd, selected))
			}

		case key.Matches(msg, m.keyMap.CompleteTask):
			if selected := m.selectedTask(); selected != nil {
				t := tasks.NewTaskFromTask(*selected, tasks.WithStatus(tasks.Done, m.ctx.Now()))
				if err := m.nd.UpdateTask(t); err != nil {
					m.footer.SetMessage(fmt.Sprintf("error completing task: %v", err), time.Now().Add(10*time.Second), m.ctx.Theme.Red)
				}
			}

		case key.Matches(msg, m.keyMap.DeleteTask):
			if selected := m.selectedTask(); selected != nil {
				if err := m.nd.DeleteTask(*selected); err != nil {
					m.footer.SetMessage(fmt.Sprintf("error deleting task: %v", err), time.Now().Add(10*time.Second), m.ctx.Theme.Red)
				}
			}
		}
	}

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
	if model != nil { // if model is not nil we're navigating to a new view
		return model, tea.Batch(command, cmd)
	}
	cmd = tea.Batch(cmd, command)
	return m, cmd

}

func (m *Model) selectedTask() *tasks.Task {
	if m.completed.Focused() {
		return m.completed.Selected()
	}
	return m.tasklist.Selected()
}

func (m *Model) togglePanels() {
	if m.tasklist.Focused() {
		m.tasklist.Blur()
		m.completed.Focus()
	} else {
		m.completed.Blur()
		m.tasklist.Focus()
	}
}

func (m *Model) moveUp(n int) {
	if m.tasklist.Focused() {
		m.tasklist.MoveUp(n)
	} else {
		m.completed.MoveUp(n)
	}
}

func (m *Model) moveDown(n int) {
	if m.tasklist.Focused() {
		m.tasklist.MoveDown(n)
	} else {
		m.completed.MoveDown(n)
	}
}

func (m *Model) updateDate(date time.Time) {
	m.date = date
	m.updateTasks()
}
