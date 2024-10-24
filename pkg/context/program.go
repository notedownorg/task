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
}

func (c *ProgramContext) Init() (tea.Model, tea.Cmd) {
	return c, nil
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
		}
	case tea.WindowSizeMsg:
		c.onWindowResize(msg)
	}
	return c, nil
}

func (c *ProgramContext) onWindowResize(msg tea.WindowSizeMsg) {
	c.ScreenHeight = msg.Height
	c.ScreenWidth = msg.Width
}
