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

package colors

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/notedownorg/task/pkg/themes"
)

func Task(theme themes.Theme, status tasks.Status) (bg, fg lipgloss.Color, err error) {
	switch status {
	case tasks.Todo:
		return theme.Panel, theme.Text, nil
	case tasks.Blocked:
		return theme.Panel, theme.Red, nil
	case tasks.Doing:
		return theme.Panel, theme.Green, nil
	default:
		return theme.Panel, theme.Text, fmt.Errorf("unexpected task status %v", status)
	}
}

func TaskSelected(theme themes.Theme, status tasks.Status) (bg, fg lipgloss.Color, err error) {
	switch status {
	case tasks.Todo:
		return theme.Text, theme.TextCursor, nil
	case tasks.Blocked:
		return theme.Red, theme.TextCursor, nil
	case tasks.Doing:
		return theme.Green, theme.TextCursor, nil
	default:
		return theme.Text, theme.TextCursor, fmt.Errorf("unexpected task status %v", status)
	}
}
