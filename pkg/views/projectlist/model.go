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

package projectlist

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/notedownorg/notedown/pkg/providers/projects"
	"github.com/notedownorg/task/pkg/components/groupedlist"
	"github.com/notedownorg/task/pkg/components/statusbar"
	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/listeners"
	"github.com/notedownorg/task/pkg/notedown"
	"github.com/notedownorg/task/pkg/views/projectadd"
	"github.com/notedownorg/task/pkg/views/projectmanager"
)

const (
	view = "projects"
)

func HandleNew(nd notedown.Client) context.GlobalKeyHandler {
	return func(ctx *context.ProgramContext, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
		key := msg.Key()
		if key.Mod == tea.ModCtrl && key.Code == 'p' {
			return ctx.Navigate(New(ctx, nd))
		}
		return nil, nil
	}
}

func New(ctx *context.ProgramContext, nd notedown.Client) *Model {
	m := &Model{
		ctx: ctx,
		nd:  nd,

		keyMap: DefaultKeyMap,

		projectlist: groupedlist.New(groupedlist.WithRenderers(mainRendererFuncs(ctx.Theme))).Focus(),
		closed:      groupedlist.New(groupedlist.WithRenderers(closedRendererFuncs(ctx.Theme))),
		footer:      statusbar.New(ctx, statusbar.NewMode(view, statusbar.ActionNeutral), nd),
	}
	m.updateProjects()
	return m
}

type Model struct {
	ctx *context.ProgramContext
	nd  notedown.Client

	keyMap KeyMap

	projectlist *groupedlist.Model[projects.Project]
	closed      *groupedlist.Model[projects.Project]
	footer      *statusbar.Model
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
		case key.Matches(msg, m.keyMap.TogglePanels):
			m.togglePanels()
		case key.Matches(msg, m.keyMap.AddProject):
			return m.ctx.Navigate(projectadd.New(m.ctx, m.nd))
		case key.Matches(msg, m.keyMap.EditProject):
			if selected := m.selectedProject(); selected != nil {
				return m.ctx.Navigate(projectmanager.New(m.ctx, m.nd, *selected))
			}
		case key.Matches(msg, m.keyMap.DeleteProject):
			if selected := m.selectedProject(); selected != nil {
				if err := m.nd.DeleteProject(*selected); err != nil {
					m.footer.SetMessage(fmt.Sprintf("error deleting project: %v", err), time.Now().Add(10*time.Second), m.ctx.Theme.Red)
				}
			}
		case key.Matches(msg, m.keyMap.CursorUp):
			m.moveUp(1)
		case key.Matches(msg, m.keyMap.CursorDown):
			m.moveDown(1)
		}
	}

	// If the task client has emitted an event, update the tasks
	if _, ok := msg.(listeners.ProjectEvent); ok {
		m.updateProjects()
	}

	// Handle program level key presses and events
	model, command := m.ctx.Update(msg)
	if model != nil { // if model is not nil we're navigating to a new view
		return model, tea.Batch(command, cmd)
	}
	cmd = tea.Batch(cmd, command)
	return m, cmd
}

func (m *Model) selectedProject() *projects.Project {
	if m.projectlist.Focused() {
		return m.projectlist.Selected()
	}
	return m.closed.Selected()
}

func (m *Model) togglePanels() {
	if m.projectlist.Focused() {
		m.projectlist.Blur()
		m.closed.Focus()
	} else {

		m.closed.Blur()
		m.projectlist.Focus()
	}
}

func (m *Model) moveUp(n int) {
	if m.projectlist.Focused() {
		m.projectlist.MoveUp(n)
	} else {
		m.closed.MoveUp(n)
	}
}

func (m *Model) moveDown(n int) {
	if m.projectlist.Focused() {
		m.projectlist.MoveDown(n)
	} else {
		m.closed.MoveDown(n)
	}
}
