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
	"fmt"
	"log/slog"

	"github.com/charmbracelet/lipgloss"
	"github.com/notedownorg/notedown/pkg/providers/projects"
	"github.com/notedownorg/task/pkg/themes"
)

var (
	s = lipgloss.NewStyle
	w = lipgloss.Width
	h = lipgloss.Height
)

func colors(theme themes.Theme, project projects.Project) (bg, fg lipgloss.Color, err error) {
	switch project.Status() {
	case projects.Backlog:
		return theme.Panel, theme.Text, nil
	case projects.Blocked:
		return theme.Panel, theme.Red, nil
	case projects.Active:
		return theme.Panel, theme.Green, nil
	default:
		return theme.Panel, theme.Text, fmt.Errorf("unexpected project status %v", project.Status())
	}
}

func selectedColors(theme themes.Theme, project projects.Project) (bg, fg lipgloss.Color, err error) {
	switch project.Status() {
	case projects.Backlog:
		return theme.Text, theme.TextCursor, nil
	case projects.Blocked:
		return theme.Red, theme.TextCursor, nil
	case projects.Active:
		return theme.Green, theme.TextCursor, nil
	default:
		return theme.Text, theme.TextCursor, fmt.Errorf("unexpected project status %v", project.Status())
	}
}

func icon(status projects.Status) string {
	switch status {
	case projects.Backlog:
		return ""
	case projects.Blocked:
		return ""
	case projects.Active:
		return ""
	case projects.Archived:
		return ""
	case projects.Abandoned:
		return ""
	default:
		slog.Warn("unknown project status", "status", status)
		return " "
	}
}
