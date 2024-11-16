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
	"log/slog"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/notedownorg/task/pkg/components/groupedlist"
	"github.com/notedownorg/task/pkg/themes"
)

func mainRendererFuncs(theme themes.Theme, dateRetriever func() time.Time) groupedlist.Renderers[tasks.Task] {
	paddingHorizontal := 2

	return groupedlist.Renderers[tasks.Task]{
		Header: func(name string, width int) string {
			bg := func(s string) lipgloss.Color {
				switch s {
				case "Doing":
					return theme.Green
				case "Todo":
					return theme.Text
				case "Blocked":
					return theme.Red
				default:
					slog.Warn("unexpected task status", "status", s)
					return theme.Text
				}
			}(name)

			return lipgloss.JoinVertical(lipgloss.Top,
				s().Margin(0, 0, 1, 0).
					Background(bg).
					Foreground(theme.TextCursor).
					Bold(true).
					Padding(0, paddingHorizontal).
					Render(strings.ToUpper(name)),
				s().Width(width).Background(theme.Panel).
					Render(""),
			)
		},

		Footer: func(name string, width int) string {
			return lipgloss.JoinVertical(lipgloss.Bottom,
				s().Width(width).Background(theme.Panel).
					Render(""),
				"",
			)
		},

		Item: func(task tasks.Task, width int) string {
			bg, fg, err := colors(theme, task)
			if err != nil {
				slog.Warn("unexpected task status", "status", task.Status())
				return ""
			}

			right := buildRight(false)(theme, task, dateRetriever, bg)

			remainingSpace := width - w(right) - 2*paddingHorizontal
			left := buildLeft(false)(theme, task, remainingSpace, bg)

			middlePadding := width - w(left) - w(right) - 2*paddingHorizontal
			middle := s().Background(bg).PaddingRight(middlePadding).Render("") // fill out the rest of the space

			return s().Width(width).Padding(0, paddingHorizontal).Background(bg).Foreground(fg).Render(left + middle + right)
		},

		Selected: func(task tasks.Task, width int) string {
			bg, fg, err := selectedColors(theme, task)
			if err != nil {
				slog.Warn("unexpected task status", "status", task.Status())
				return ""
			}

			right := buildRight(true)(theme, task, dateRetriever, bg)

			remainingSpace := width - w(right) - 2*paddingHorizontal
			left := buildLeft(true)(theme, task, remainingSpace, bg)

			middlePadding := width - w(left) - w(right) - 2*paddingHorizontal
			middle := s().Background(bg).PaddingRight(middlePadding).Render("") // fill out the rest of the space

			return s().Width(width).Padding(0, paddingHorizontal).Background(bg).Foreground(fg).Render(left + middle + right)
		},
	}
}

// See https://github.com/charmbracelet/lipgloss/issues/144 for why we need to pass bg
func buildRight(selected bool) func(theme themes.Theme, task tasks.Task, dateRetriever func() time.Time, bg lipgloss.Color) string {
	return func(theme themes.Theme, task tasks.Task, dateRetriever func() time.Time, bg lipgloss.Color) string {
		res := make([]string, 0)

		// Prefer due date over scheduled date as it is definitially more important
		if task.Due() != nil {
			due, date := *task.Due(), dateRetriever().Truncate(time.Hour*24)
			due = due.Truncate(time.Hour * 24)
			dr := ""
			if dateBefore(due, date) {
				if (date.Sub(due) / (24 * time.Hour)) == 1 {
					dr = "󰃭 Yesterday"
				} else {
					dr = "󰃭 " + task.Due().Format("Jan  _2"+dayOfMonthSuffix(due.Day()))
				}
			}
			if selected {
				res = append(res, s().Background(bg).Foreground(theme.TextCursor).Render(dr))
			} else {
				res = append(res, s().Background(bg).Foreground(theme.Red).Render(dr))
			}
		} else if task.Scheduled() != nil {
			scheduled, date := *task.Scheduled(), dateRetriever().Truncate(time.Hour*24)
			scheduled = scheduled.Truncate(time.Hour * 24)
			sr := ""
			if dateBefore(scheduled, date) {
				if (date.Sub(scheduled) / (24 * time.Hour)) == 1 {
					sr = "󰃭 Yesterday"
				} else {
					sr = "󰃭 " + task.Scheduled().Format("Jan  _2"+dayOfMonthSuffix(scheduled.Day()))
				}
			}
			if selected {
				res = append(res, s().Background(bg).Foreground(theme.TextCursor).Render(sr))
			} else {
				res = append(res, s().Background(bg).Foreground(theme.Text).Render(sr))
			}
		}

		return strings.Join(res, "  ")
	}
}

// See https://github.com/charmbracelet/lipgloss/issues/144 for why we need to pass bg
func buildLeft(selected bool) func(theme themes.Theme, task tasks.Task, remainingSpace int, bg lipgloss.Color) string {
	return func(theme themes.Theme, task tasks.Task, remainingSpace int, bg lipgloss.Color) string {
		res := make([]string, 0)

		i := icon(task.Status())
		res = append(res, i, "  ")

		e := every(task)

		textWidth := remainingSpace - w(i) - w(e)
		text := runewidth.Truncate(task.Name(), textWidth, "…")
		res = append(res, text)

		if e != "" {
			if selected {
				res = append(res, " ", s().Background(bg).Foreground(theme.TextCursor).Render(e))
			} else {
				res = append(res, " ", s().Background(bg).Foreground(theme.Text).Render(e))
			}
		}
		return strings.Join(res, "")
	}
}
