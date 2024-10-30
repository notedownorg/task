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

func renderTask(theme themes.Theme, task ast.Task, selected bool) string {
	i := icon(task.Status())
	name := task.Name()

	// Handle styling based on task status
	s := lipgloss.NewStyle()
	switch task.Status() {
	case ast.Doing:
		if selected {
			return s.Padding(0, 1).Background(theme.Green).Foreground(theme.TextCursor).Render(i + "  " + name)
		}
		return s.Padding(0, 1).Background(theme.Panel).Foreground(theme.Green).Render(i + "  " + name)
	case ast.Todo:
		if selected {
			return s.Padding(0, 1).Background(theme.Text).Foreground(theme.TextCursor).Render(i + "  " + name)
		}
		return s.Padding(0, 1).Background(theme.Panel).Foreground(theme.Text).Render(i + "  " + name)
	case ast.Blocked:
		if selected {
			return s.Padding(0, 1).Background(theme.Red).Foreground(theme.TextCursor).Render(i + "  " + name)
		}
		return s.Padding(0, 1).Background(theme.Panel).Foreground(theme.Red).Render(i + "  " + name)
	case ast.Done, ast.Abandoned:
		if selected {
			base := s.Background(theme.TextFaint).Foreground(theme.TextCursor)
			return base.Padding(0, 1).Render(base.Render(i+"  ") + base.Strikethrough(true).Render(name))
		}
		base := s.Background(theme.Panel).Foreground(theme.TextFaint)
		return s.Background(theme.Panel).Padding(0, 1).Render(base.Render(i+"  ") + base.Strikethrough(true).Render(name))
	}
	slog.Warn("unknown task status", "status", task.Status())
	return ""
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
