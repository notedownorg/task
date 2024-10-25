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

package tasklist

import (
	"log/slog"

	"github.com/charmbracelet/lipgloss"
	"github.com/notedownorg/notedown/pkg/ast"
	"github.com/notedownorg/task/pkg/themes"
)

func taskStyle(theme themes.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.Text)
}

func selectedTaskStyle(theme themes.Theme) lipgloss.Style {
	return lipgloss.NewStyle().Inherit(taskStyle(theme)).
		Foreground(theme.TextCursor).
		Background(theme.Text).
		Bold(true)
}

func (m Model) renderTask(task ast.Task, selected bool) string {
	icon := m.renderIcon(task.Status())
	space := " "
	name := task.Name()

	t := lipgloss.JoinHorizontal(lipgloss.Left, icon, space, space, name)

	if selected {
		return selectedTaskStyle(m.ctx.Theme).Render(t)
	}
	return taskStyle(m.ctx.Theme).Render(t)
}

func (m Model) renderIcon(status ast.Status) string {
	render := lipgloss.NewStyle().Padding(0, 1).Render
	switch status {
	case ast.Todo:
		return render("")
	case ast.Blocked:
		return render("")
	case ast.Doing:
		return render("") //  or  ?
	case ast.Done:
		return render("")
	case ast.Abandoned:
		return render("")
	default:
		slog.Warn("unknown task status", "status", status)
		return render(" ")
	}
}
