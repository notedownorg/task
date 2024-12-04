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
	"log/slog"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/notedownorg/notedown/pkg/providers/projects"
	"github.com/notedownorg/task/pkg/components/groupedlist"
	"github.com/notedownorg/task/pkg/styling/colors"
	"github.com/notedownorg/task/pkg/styling/icons"
	"github.com/notedownorg/task/pkg/themes"
)

func mainRendererFuncs(theme themes.Theme) groupedlist.Renderers[projects.Project] {
	paddingHorizontal := 2

	return groupedlist.Renderers[projects.Project]{
		Header: func(name string, width int) string {
			bg := func(s string) lipgloss.Color {
				switch s {
				case "Active":
					return theme.Green
				case "Backlog":
					return theme.Text
				case "Blocked":
					return theme.Red
				default:
					slog.Warn("unexpected project status", "status", s)
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

		Item: func(project projects.Project, width int) string {
			bg, fg, err := colors.Project(theme, project.Status())
			if err != nil {
				slog.Warn("unexpected project status", "status", project.Status())
				return ""
			}

			right := buildRight(false)(theme, project, bg)

			remainingSpace := width - w(right) - 2*paddingHorizontal
			left := buildLeft(false)(theme, project, remainingSpace, bg)

			middlePadding := width - w(left) - w(right) - 2*paddingHorizontal
			middle := s().Background(bg).PaddingRight(middlePadding).Render("") // fill out the rest of the space

			return s().Width(width).Padding(0, paddingHorizontal).Background(bg).Foreground(fg).Render(left + middle + right)
		},

		Selected: func(project projects.Project, width int) string {
			bg, fg, err := colors.ProjectSelected(theme, project.Status())
			if err != nil {
				slog.Warn("unexpected project status", "status", project.Status())
				return ""
			}

			right := buildRight(true)(theme, project, bg)

			remainingSpace := width - w(right) - 2*paddingHorizontal
			left := buildLeft(true)(theme, project, remainingSpace, bg)

			middlePadding := width - w(left) - w(right) - 2*paddingHorizontal
			middle := s().Background(bg).PaddingRight(middlePadding).Render("") // fill out the rest of the space

			return s().Width(width).Padding(0, paddingHorizontal).Background(bg).Foreground(fg).Render(left + middle + right)
		},
	}
}

// See https://github.com/charmbracelet/lipgloss/issues/144 for why we need to pass bg
func buildRight(selected bool) func(theme themes.Theme, project projects.Project, bg lipgloss.Color) string {
	return func(theme themes.Theme, project projects.Project, bg lipgloss.Color) string {
		return ""
	}
}

// See https://github.com/charmbracelet/lipgloss/issues/144 for why we need to pass bg
func buildLeft(selected bool) func(theme themes.Theme, project projects.Project, remainingSpace int, bg lipgloss.Color) string {
	return func(theme themes.Theme, project projects.Project, remainingSpace int, bg lipgloss.Color) string {
		res := make([]string, 0)

		i := icons.Project(project.Status())
		res = append(res, i, "  ")

		textWidth := remainingSpace - w(i)
		text := runewidth.Truncate(project.Name(), textWidth, "…")
		res = append(res, text)

		return strings.Join(res, "")
	}
}
