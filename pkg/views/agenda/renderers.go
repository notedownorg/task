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

	"github.com/charmbracelet/lipgloss"
	"github.com/notedownorg/notedown/pkg/ast"
	"github.com/notedownorg/task/pkg/components/groupedlist"
	"github.com/notedownorg/task/pkg/themes"
)

func rendererFuncs(theme themes.Theme) groupedlist.Renderers[ast.Task] {
	s := lipgloss.NewStyle()
	paddingHorizontal := 2

	return groupedlist.Renderers[ast.Task]{
		Header: func(name string, width int) string {
			return lipgloss.JoinVertical(lipgloss.Top,
				s.Margin(0, 0, 1, 0).
					Render(name),
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
			i := icon(task.Status())
			name := task.Name()

			switch task.Status() {
			case ast.Doing:
				return s.Width(width).Padding(0, paddingHorizontal).Background(theme.Panel).Foreground(theme.Green).Render(i + "  " + name)
			case ast.Todo:
				return s.Width(width).Padding(0, paddingHorizontal).Background(theme.Panel).Foreground(theme.Text).Render(i + "  " + name)
			case ast.Blocked:
				return s.Width(width).Padding(0, paddingHorizontal).Background(theme.Panel).Foreground(theme.Red).Render(i + "  " + name)
			case ast.Done, ast.Abandoned:
				base := s.Background(theme.Panel).Foreground(theme.TextFaint)
				return s.Width(width).Background(theme.Panel).Padding(0, paddingHorizontal).Render(base.Render(i+"  ") + base.Strikethrough(true).Render(name))
			}
			slog.Warn("unknown task status", "status", task.Status())
			return ""
		},

		Selected: func(task ast.Task, width int) string {
			i := icon(task.Status())
			name := task.Name()

			switch task.Status() {
			case ast.Doing:
				return s.Width(width).Padding(0, paddingHorizontal).Background(theme.Green).Foreground(theme.TextCursor).Render(i + "  " + name)
			case ast.Todo:
				return s.Width(width).Padding(0, paddingHorizontal).Background(theme.Text).Foreground(theme.TextCursor).Render(i + "  " + name)
			case ast.Blocked:
				return s.Width(width).Padding(0, paddingHorizontal).Background(theme.Red).Foreground(theme.TextCursor).Render(i + "  " + name)
			case ast.Done, ast.Abandoned:
				base := s.Background(theme.TextFaint).Foreground(theme.TextCursor)
				return s.Width(width).Padding(0, paddingHorizontal).Background(theme.TextFaint).Render(base.Render(i+"  ") + base.Strikethrough(true).Render(name))
			}

			slog.Warn("unknown task status", "status", task.Status())
			return ""
		},
	}
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
