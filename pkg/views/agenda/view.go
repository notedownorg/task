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
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	s = lipgloss.NewStyle
	w = lipgloss.Width
	h = lipgloss.Height
)

func (m *Model) View() string {
	gap := "   "
	horizontalPadding := 2
	verticalPadding := 1

	header := fmt.Sprintf("← %v →", humanizeDate(m.date, time.Now()))

	footer := m.footer.
		Width(m.ctx.ScreenWidth - horizontalPadding*2).
		View()

	completed := m.completed.
		Height(m.ctx.ScreenHeight - h(footer) - h(header) - verticalPadding*2 - 2). // -2 for the gaps
		Width(m.ctx.ScreenWidth/4 - horizontalPadding*2).
		View()

	tasklist := m.tasklist.
		Height(m.ctx.ScreenHeight - h(footer) - h(header) - verticalPadding*2 - 2). // -2 for the gaps
		Width(m.ctx.ScreenWidth - w(completed) - horizontalPadding*2 - 3).          // -3 for the gap
		View()

	main := lipgloss.JoinHorizontal(lipgloss.Left, tasklist, gap, completed)
	panel := lipgloss.JoinVertical(lipgloss.Top, header, gap, main, gap, footer)

	return lipgloss.NewStyle().Padding(verticalPadding, horizontalPadding).Render(panel)
}

func humanizeDate(date time.Time, relativeTo time.Time) string {
	rel := time.Date(relativeTo.Year(), relativeTo.Month(), relativeTo.Day(), 0, 0, 0, 0, time.UTC)
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	diff := date.Sub(rel) / (24 * time.Hour)

	switch diff {
	case -6, -5, -4, -3, -2:
		return fmt.Sprintf("Last %s", date.Weekday().String())
	case -1:
		return "Yesterday"
	case 0:
		return "Today"
	case 1:
		return "Tomorrow"
	case 2, 3, 4, 5, 6:
		return date.Weekday().String()
	}

	if date.Year() == relativeTo.Year() {
		return date.Format("Monday, January 2")
	}
	return date.Format("Monday, January 2, 2006")
}
