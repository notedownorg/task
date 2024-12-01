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

package taskreschedule

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/notedownorg/task/pkg/components/statusbar"
	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/model"
	"github.com/notedownorg/task/pkg/notedown"
)

type Model struct {
	base model.Base
	ctx  *context.ProgramContext
	nd   notedown.Client

	original *tasks.Task
	date     time.Time

	keyMap KeyMap

	footer *statusbar.Model
}

func New(ctx *context.ProgramContext, nd notedown.Client, task *tasks.Task) *Model {
	date := time.Date(ctx.Now().Year(), ctx.Now().Month(), ctx.Now().Day(), 0, 0, 0, 0, ctx.Now().Location())
	m := &Model{
		ctx:      ctx,
		nd:       nd,
		original: task,
		date:     date,

		keyMap: DefaultKeyMap,
		footer: statusbar.New(ctx, statusbar.NewMode("reschedule task", statusbar.ActionEdit), nd),
	}
	return m
}

func (m *Model) Init() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle view level key presses
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Today):
			return m.submit(m.date)
		case key.Matches(msg, m.keyMap.Tommorrow):
			return m.submit(m.date.AddDate(0, 0, 1))
		case key.Matches(msg, m.keyMap.InTwoDays):
			return m.submit(m.date.AddDate(0, 0, 2))
		case key.Matches(msg, m.keyMap.InThreeDays):
			return m.submit(m.date.AddDate(0, 0, 3))
		case key.Matches(msg, m.keyMap.InFourDays):
			return m.submit(m.date.AddDate(0, 0, 4))
		case key.Matches(msg, m.keyMap.InFiveDays):
			return m.submit(m.date.AddDate(0, 0, 5))
		case key.Matches(msg, m.keyMap.InSixDays):
			return m.submit(m.date.AddDate(0, 0, 6))
		case key.Matches(msg, m.keyMap.InSevenDays):
			return m.submit(m.date.AddDate(0, 0, 7))
		case key.Matches(msg, m.keyMap.InFortnight):
			return m.submit(m.date.AddDate(0, 0, 14))
		case key.Matches(msg, m.keyMap.NextMonth):
			return m.submit(time.Date(m.date.Year(), m.date.Month()+1, 1, 0, 0, 0, 0, m.date.Location()))
		case key.Matches(msg, m.keyMap.NextYear):
			return m.submit(time.Date(m.date.Year()+1, 1, 1, 0, 0, 0, 0, m.date.Location()))
		}
	}

	// Handle program level key presses and events
	model, command := m.ctx.Update(msg)
	if model != nil { // if model is not nil we're navigating to a new view
		return model, tea.Batch(command, cmd)
	}
	cmd = tea.Batch(cmd, command)
	return m, cmd
}

func (m *Model) View() string {
	horizontalPadding := 2
	verticalMargin := 1

	footer := m.footer.
		Width(m.ctx.ScreenWidth-horizontalPadding*2).
		Margin(verticalMargin, 0).
		View()

	optRender := func(date time.Time, key string, helper string) string {
		suffix := lipgloss.NewStyle().Foreground(m.ctx.Theme.TextFaint).Render(fmt.Sprintf("[%s]", helper))
		return m.base.NewStyle().Render(
			fmt.Sprintf("%s 󰁕 %s ", key, date.Format("2006-01-02")) + suffix,
		)
	}

	top := lipgloss.NewStyle().
		Margin(1, 3).
		Render(
			lipgloss.JoinVertical(lipgloss.Top,
				optRender(m.date, "0", "today"),
				optRender(m.date.AddDate(0, 0, 1), "1", "tomorrow"),
				optRender(m.date.AddDate(0, 0, 2), "2", "in two days"),
				optRender(m.date.AddDate(0, 0, 3), "3", "in three days"),
				optRender(m.date.AddDate(0, 0, 4), "4", "in four days"),
				optRender(m.date.AddDate(0, 0, 5), "5", "in five days"),
				optRender(m.date.AddDate(0, 0, 6), "6", "in six days"),
				optRender(m.date.AddDate(0, 0, 7), "7", "in seven days"),
				optRender(m.date.AddDate(0, 0, 14), "f", "in a fortnight"),
				optRender(time.Date(m.date.Year(), m.date.Month()+1, 1, 0, 0, 0, 0, m.date.Location()), "m", "next month"),
				optRender(time.Date(m.date.Year()+1, 1, 1, 0, 0, 0, 0, m.date.Location()), "y", "next year"),
			),
		)

	border := lipgloss.RoundedBorder()
	var b strings.Builder
	str := "Reschedule-Task"
	for i := len(str) + 2; i <= lipgloss.Width(top); i++ {
		b.WriteString(lipgloss.RoundedBorder().Top)
	}
	b.WriteString(str)
	border.Top = b.String()

	color := m.ctx.Theme.Yellow
	form := lipgloss.NewStyle().
		Border(border).
		BorderForeground(color).
		Render(top)

	width := m.ctx.ScreenWidth - horizontalPadding*2
	height := m.ctx.ScreenHeight - lipgloss.Height(footer)

	dialog := lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, form)

	panel := lipgloss.JoinVertical(lipgloss.Top, dialog, footer)

	return lipgloss.NewStyle().Padding(0, horizontalPadding).Render(panel)
}
