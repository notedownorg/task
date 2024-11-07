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

package taskeditor

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/notedownorg/notedown/pkg/fileserver/writer"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
)

func (m *Model) submit() (tea.Model, tea.Cmd) {
	// Build the task
	opts := make([]tasks.TaskOption, 0)
	if m.fields.Due != nil {
		opts = append(opts, tasks.WithDue(*m.fields.Due))
	}
	if m.fields.Scheduled != nil {
		opts = append(opts, tasks.WithScheduled(*m.fields.Scheduled))
	}
	if m.fields.Priority != nil {
		opts = append(opts, tasks.WithPriority(*m.fields.Priority))
	}
	if m.fields.Every != nil {
		opts = append(opts, tasks.WithEvery(*m.fields.Every))
	}
	if m.fields.Completed != nil {
		opts = append(opts, tasks.WithCompleted(*m.fields.Completed))
	}

	// Create/Update are intentionally run syncronously to prevent losing progress on error
	if m.mode == adding {
		if err := m.tasks.Create(m.location.file, writer.AT_END, m.fields.Name, m.status.Value(), opts...); err != nil {
			slog.Error("failed to create task", "error", err)

			// TODO: We should probably show an error message to the user
			return m, nil
		}

		// If we've successfully created the task, we can navigate back to the previous view
		return m.ctx.Back(), nil
	}

	// Ensure name/status are also updated
	opts = append(opts, tasks.WithName(m.fields.Name))
	opts = append(opts, tasks.WithStatus(m.status.Value()))

	task := tasks.NewTaskFromTask(*m.original, opts...)
	if err := m.tasks.Update(task); err != nil {
		slog.Error("failed to update task", "error", err)

		// TODO: We should probably show an error message to the user
		return m, nil
	}

	// If we've successfully created the task, we can navigate back to the previous view
	return m.ctx.Back(), nil
}
