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
	"log/slog"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/notedownorg/notedown/pkg/ast"
	"github.com/notedownorg/task/pkg/components/groupedlist"
	"github.com/notedownorg/task/pkg/themes"
)

func mainRendererFuncs(theme themes.Theme) groupedlist.Renderers[ast.Task] {
	s := lipgloss.NewStyle()
	paddingHorizontal := 2

	return groupedlist.Renderers[ast.Task]{
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
				s.Margin(0, 0, 1, 0).
					Background(bg).
					Foreground(theme.TextCursor).
					Bold(true).
					Padding(0, paddingHorizontal).
					Render(strings.ToUpper(name)),
				s.Width(width).Background(theme.Panel).
					Render(""),
			)
		},

		Footer: func(name string, width int) string {
			return lipgloss.JoinVertical(lipgloss.Bottom,
				s.Width(width).Background(theme.Panel).
					Render(""),
				"",
			)
		},

		Item: func(task ast.Task, width int) string {
			fields := []string{
				icon(task.Status()),
				lipgloss.NewStyle().Width(70).Render(runewidth.Truncate(task.Name(), 70, "…")),
				priority(task),
			}

			switch task.Status() {
			case ast.Doing:
				return s.Width(width).Padding(0, paddingHorizontal).Background(theme.Panel).Foreground(theme.Green).Render(strings.Join(fields, "  "))
			case ast.Todo:
				return s.Width(width).Padding(0, paddingHorizontal).Background(theme.Panel).Foreground(theme.Text).Render(strings.Join(fields, "  "))
			case ast.Blocked:
				return s.Width(width).Padding(0, paddingHorizontal).Background(theme.Panel).Foreground(theme.Red).Render(strings.Join(fields, "  "))
			}
			slog.Warn("unknown task status", "status", task.Status())
			return ""
		},

		Selected: func(task ast.Task, width int) string {
			fields := []string{
				icon(task.Status()),
				lipgloss.NewStyle().Width(70).Render(runewidth.Truncate(task.Name(), 70, "…")),
				priority(task),
			}

			switch task.Status() {
			case ast.Doing:
				return s.Width(width).Padding(0, paddingHorizontal).Background(theme.Green).Foreground(theme.TextCursor).Render(strings.Join(fields, "  "))
			case ast.Todo:
				return s.Width(width).Padding(0, paddingHorizontal).Background(theme.Text).Foreground(theme.TextCursor).Render(strings.Join(fields, "  "))
			case ast.Blocked:
				return s.Width(width).Padding(0, paddingHorizontal).Background(theme.Red).Foreground(theme.TextCursor).Render(strings.Join(fields, "  "))
			case ast.Done, ast.Abandoned:
				base := s.Background(theme.TextFaint).Foreground(theme.TextCursor)
				return s.Width(width).Padding(0, paddingHorizontal).Background(theme.TextFaint).Render(base.Render(fields[0]+"  ") + base.Strikethrough(true).Render(fields[1]))
			}

			slog.Warn("unknown task status", "status", task.Status())
			return ""
		},
	}
}

func completedRendererFuncs(theme themes.Theme) groupedlist.Renderers[ast.Task] {
	s := lipgloss.NewStyle()
	paddingHorizontal := 2

	return groupedlist.Renderers[ast.Task]{
		Header: func(name string, width int) string {
			return lipgloss.JoinVertical(lipgloss.Top,
				s.Margin(0, 0, 1, 0).
					Background(theme.TextFaint).
					Foreground(theme.TextCursor).
					Bold(true).
					Padding(0, paddingHorizontal).
					Render(strings.ToUpper(name)),
				s.Width(width).Background(theme.Panel).
					Render(""),
			)
		},
		Footer: func(name string, width int) string {
			return lipgloss.JoinVertical(lipgloss.Bottom,
				s.Width(width).Background(theme.Panel).
					Render(""),
				"",
			)
		},
		Item: func(task ast.Task, width int) string {
			fields := parts(task, width-8)

			switch task.Status() {
			case ast.Done, ast.Abandoned:
				base := s.Background(theme.Panel).Foreground(theme.TextFaint)
				return s.Width(width).Background(theme.Panel).Padding(0, paddingHorizontal).Render(base.Render(fields[0]+"  ") + base.Strikethrough(true).Render(fields[1]))
			}

			slog.Warn("unexpected task status", "status", task.Status())
			return ""
		},
		Selected: func(task ast.Task, width int) string {
			fields := parts(task, width-8)

			switch task.Status() {
			case ast.Done, ast.Abandoned:
				base := s.Background(theme.TextFaint).Foreground(theme.TextCursor)
				return s.Width(width).Padding(0, paddingHorizontal).Background(theme.TextFaint).Render(base.Render(fields[0]+"  ") + base.Strikethrough(true).Render(fields[1]))
			}

			slog.Warn("unexpected task status", "status", task.Status())
			return ""
		},
	}
}

func parts(task ast.Task, trunc int) []string {
	res := []string{
		icon(task.Status()),
		lipgloss.NewStyle().Width(trunc).Render(runewidth.Truncate(task.Name(), trunc, "…")),
	}
	if task.Priority() != nil {
		res = append(res, fmt.Sprintf(" %d", *task.Priority()))
	} else {
		res = append(res, "   ")
	}
	return res
}

func priority(task ast.Task) string {
	if task.Priority() != nil {
		return fmt.Sprintf(" %d", *task.Priority())
	}
	return "   "
}

func icon(status ast.Status) string {
	switch status {
	case ast.Todo:
		return ""
	case ast.Blocked:
		return ""
	case ast.Doing:
		return ""
	case ast.Done:
		return ""
	case ast.Abandoned:
		return ""
	default:
		slog.Warn("unknown task status", "status", status)
		return " "
	}
}
