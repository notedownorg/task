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
	"time"

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

	KeyHandlers []GlobalKeyHandler

	// clock acts as the "system" clock for the program, if nil uses time.Now()
	// typically this would only be set (pinned) for testing purposes.
	clock func() time.Time

	initialView tea.Model
}

type ProgramContextOption func(*ProgramContext)

func WithClock(fn func() time.Time) ProgramContextOption {
	return func(p *ProgramContext) {
		p.clock = fn
	}
}

type InitalViewBuilder func(*ProgramContext) tea.Model

func New(theme themes.Theme, initial InitalViewBuilder, opts ...ProgramContextOption) *ProgramContext {
	p := &ProgramContext{Theme: theme}

	for _, opt := range opts {
		opt(p)
	}

	// Build the initial view after the options have been applied in case they affect the view.
	p.initialView, _ = p.Navigate(initial(p))
	return p
}

func (c *ProgramContext) Init() (tea.Model, tea.Cmd) {
	return c.initialView, c.Listeners.Init()
}

func (c *ProgramContext) View() string {
	return ""
}

// Note you probably want to call this at the end of the update function to allow you to override the global key handlers.
func (c *ProgramContext) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, handler := range c.KeyHandlers {
			if m, cmd := handler(c, msg); m != nil || cmd != nil {
				return m, cmd
			}
		}
	case tea.WindowSizeMsg:
		c.onWindowResize(msg)
	}
	return nil, c.Listeners.Receive(msg)
}

func (c ProgramContext) Now() time.Time {
	if c.clock == nil {
		return time.Now()
	}
	return c.clock()
}

func (c *ProgramContext) onWindowResize(msg tea.WindowSizeMsg) {
	c.ScreenHeight = msg.Height
	c.ScreenWidth = msg.Width
}
