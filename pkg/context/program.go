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

package context

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/notedownorg/task/pkg/themes"
)

// ProgramContext is responsible for managing "global" state and actions for the program.
// This includes things like the theme, screen dimensions and exiting.
type ProgramContext struct {
	Theme themes.Theme

	ScreenHeight int
	ScreenWidth  int

	// History is a stack of views that the user has navigated through.
	History History

	// Listeners are a group of subscribers that will be sent every message and optionally respond with a command.
	// Examples might include a timer or an event listener that trigger re-renders.
	// tea.Model Update is intentionally not used here as we don't want listeners to be able to change the model.
	Listeners Listeners
}

type ProgramContextOption func(*ProgramContext)

func New(theme themes.Theme, opts ...ProgramContextOption) *ProgramContext {
	p := &ProgramContext{Theme: theme}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (c *ProgramContext) Init() (tea.Model, tea.Cmd) {
	return c, c.Listeners.Init()
}

func (c *ProgramContext) View() string {
	return ""
}

func (c *ProgramContext) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return c, tea.Quit
		case "esc":
			return c.Back(), nil
		}
	case tea.WindowSizeMsg:
		c.onWindowResize(msg)
	}
	return nil, c.Listeners.Receive(msg)
}

func (c *ProgramContext) Back() tea.Model {
	if m, ok := c.History.Pop(); ok {
		return m
	}
	// Return nil if there is no history to go back to.
	return nil
}

func (c *ProgramContext) Navigate(curr, next tea.Model) (tea.Model, tea.Cmd) {
	c.History.Push(curr)
	return next, nil
}

func (c *ProgramContext) onWindowResize(msg tea.WindowSizeMsg) {
	c.ScreenHeight = msg.Height
	c.ScreenWidth = msg.Width
}
