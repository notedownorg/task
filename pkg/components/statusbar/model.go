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

package statusbar

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/notedownorg/notedown/pkg/workspace/tasks"
	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/model"
	"github.com/notedownorg/task/pkg/themes"
)

func modeStyle(mode Mode) func(themes.Theme) lipgloss.Style {
	return func(theme themes.Theme) lipgloss.Style {
		background := func() lipgloss.Color {
			switch mode.action {
			case ActionCreate:
				return theme.Green
			case ActionEdit:
				return theme.Yellow
			case ActionDelete:
				return theme.Red
			default:
				return theme.Blue
			}
		}()
		return lipgloss.NewStyle().
			Foreground(theme.TextCursor).
			Background(background).
			Bold(true).
			Padding(0, 1)
	}
}

func textStyle(theme themes.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(theme.Panel)
}

func statsStyle(theme themes.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.TextCursor).
		Background(theme.Blue).
		Padding(0, 1)
}

type Model struct {
	base model.Base

	ctx   *context.ProgramContext
	tasks *tasks.Client

	mode Mode
}

func New(ctx *context.ProgramContext, mode Mode, t *tasks.Client) *Model {
	return &Model{
		ctx:   ctx,
		mode:  mode,
		tasks: t,
	}
}

func (m *Model) Init() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) Width(width int) *Model {
	m.base.Width(width)
	return m
}

func (m *Model) Margin(margin ...int) *Model {
	m.base.Margin(margin...)
	return m
}

func (m *Model) View() string {
	stats := fmt.Sprintf("󰄬 %d 󰢨 %d 󰧮 %d",
		len(m.tasks.ListTasks(tasks.FetchAllTasks())),
		len(m.tasks.ListDocuments(tasks.FetchAllDocuments(), tasks.FilterByDocumentType("project"))),
		len(m.tasks.ListDocuments(tasks.FetchAllDocuments())),
	)

	statsBlock := statsStyle(m.ctx.Theme).Render(stats)
	modeBlock := modeStyle(m.mode)(m.ctx.Theme).Render(strings.ToUpper(m.mode.text))

	w := lipgloss.Width
	statusBlockWidth := m.base.AvailableWidth() - w(statsBlock) - w(modeBlock)
	statusBlock := textStyle(m.ctx.Theme).Width(statusBlockWidth).Render("")

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		modeBlock,
		statusBlock,
		statsBlock,
	)

	return m.base.NewStyle().
		Render(bar)
}
