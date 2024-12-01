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

package taskreschedule

import (
	"log/slog"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
)

func (m *Model) submit(date time.Time) (tea.Model, tea.Cmd) {
	opts := make([]tasks.TaskOption, 0)

	// If scheduled date is set update it
	if m.original.Scheduled() != nil {
		opts = append(opts, tasks.WithScheduled(date))
	}

	// If due date is set update it
	if m.original.Due() != nil {
		opts = append(opts, tasks.WithDue(date))
	}

	task := tasks.NewTaskFromTask(*m.original, opts...)
	slog.Debug("submitting rescheduled task", "identifier", task.Identifier().String(), "task", task.String())
	if err := m.nd.UpdateTask(task); err != nil {
		slog.Error("failed to update task", "error", err)

		// TODO: We should probably show an error message to the user
		return m, nil
	}

	// If we've successfully created the task, we can navigate back to the previous view
	return m.ctx.Back(), nil
}
