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

package taskeditor

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/notedownorg/task/pkg/components/statusbar"
	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/notedown"
)

type mode int

const (
	adding mode = iota
	editing
)

type Model struct {
	ctx *context.ProgramContext
	nd  notedown.Client

	mode     mode
	original *tasks.Task
	date     time.Time

	keyMap KeyMap

	status   *Status
	text     *Text
	fields   *Fields
	location *Location

	footer *statusbar.Model
}

func New(ctx *context.ProgramContext, nd notedown.Client, mode Mode) *Model {
	m := &Model{
		ctx: ctx,
		nd:  nd,

		keyMap: DefaultKeyMap,
	}
	mode(m)
	return m
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.ToggleFocus):
			m.toggleFocus()
		case key.Matches(msg, m.keyMap.Submit):
			return m.submit()
		}
	}

	// Handle component events
	m.status.Update(msg)
	m.text.Update(msg)

	// Attempt to parse the full task and use the response to update the fields subcomponent
	m.parseTask()

	return m, cmd
}

func (m *Model) View() string {
	horizontalPadding := 2
	verticalMargin := 1
	w := lipgloss.Width

	footer := m.footer.
		Width(m.ctx.ScreenWidth-horizontalPadding*2).
		Margin(verticalMargin, 0).
		View()

	status := m.status.Margin(0, 2, 0, 0).View()
	text := m.text.Width(80).View()

	top := lipgloss.NewStyle().
		Margin(1, 3).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, status, text))

	fields := lipgloss.NewStyle().
		Margin(0, 3, 1, 3).
		PaddingTop(1).
		Width(w(status)+w(text)).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(m.ctx.Theme.BorderFaint).
		Render(m.fields.View())

	location := lipgloss.NewStyle().
		Margin(0, 3, 1, 3).
		Foreground(m.ctx.Theme.TextFaint).
		Render(m.location.View())

	lines := lipgloss.JoinVertical(lipgloss.Top,
		top,
		fields,
		location,
	)

	border := lipgloss.RoundedBorder()
	var b strings.Builder
	str := "Add-Task"
	if m.mode == editing {
		str = "Edit-Task"
	}
	for i := len(str) + 2; i <= lipgloss.Width(lines); i++ {
		b.WriteString(lipgloss.RoundedBorder().Top)
	}
	b.WriteString(str)
	border.Top = b.String()

	color := m.ctx.Theme.Green
	if m.mode == editing {
		color = m.ctx.Theme.Yellow
	}
	form := lipgloss.NewStyle().
		Border(border).
		BorderForeground(color).
		Render(lines)

	width := m.ctx.ScreenWidth - horizontalPadding*2
	height := m.ctx.ScreenHeight - lipgloss.Height(footer)

	dialog := lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, form)

	panel := lipgloss.JoinVertical(lipgloss.Top, dialog, footer)

	return lipgloss.NewStyle().Padding(0, horizontalPadding).Render(panel)
}

func (m *Model) toggleFocus() {
	if m.status.focused {
		m.status = m.status.Blur()
		m.text = m.text.Focus()
	} else {
		m.status = m.status.Focus()
		m.text = m.text.Blur()
	}
}
