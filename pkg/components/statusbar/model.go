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

	"github.com/charmbracelet/lipgloss"
	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/themes"
)

func viewStyle(theme themes.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(theme.TextCursor).
		Background(theme.Green).
		Bold(true).
		Padding(0, 1)
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

type Stats struct {
	Tasks     int
	Projects  int
	Documents int
}

func (s Stats) String() string {
	return fmt.Sprintf("󰄬 %d 󰢨 %d 󰧮 %d", s.Tasks, s.Projects, s.Documents)
}

func Render(ctx *context.ProgramContext, view string, stats Stats) string {
	margin := 2

	statsBlock := statsStyle(ctx.Theme).Render(stats.String())
	viewBlock := viewStyle(ctx.Theme).Render(strings.ToUpper(view))

	w := lipgloss.Width
	statusBlockWidth := ctx.ScreenWidth - w(statsBlock) - w(viewBlock) - margin*2
	statusBlock := textStyle(ctx.Theme).Width(statusBlockWidth).Render("")

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		viewBlock,
		statusBlock,
		statsBlock,
	)

	return lipgloss.NewStyle().
		Margin(0, margin).
		Width(ctx.ScreenWidth).
		Render(bar)
}
