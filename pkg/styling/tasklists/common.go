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

package tasklists

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
)

var (
	s = lipgloss.NewStyle
	w = lipgloss.Width
	h = lipgloss.Height
)

func dateBefore(a, b time.Time) bool {
	return a.Truncate(time.Hour * 24).Before(b.Truncate(time.Hour * 24))
}

func priority(task tasks.Task) string {
	if task.Priority() != nil {
		return fmt.Sprintf(" %d", *task.Priority())
	}
	return ""
}

func every(task tasks.Task) string {
	if task.Every() != nil {
		return "󰕇"
	}
	return ""
}

func dayOfMonthSuffix(day int) string {
	switch day {
	case 1, 21, 31:
		return "st"
	case 2, 22:
		return "nd"
	case 3, 23:
		return "rd"
	default:
		return "th"
	}
}
