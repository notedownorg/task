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
	"fmt"
	"strings"
	"time"

	"github.com/a-h/parse"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
)

func (m *Model) parseTask() {
	// Build it into a task string and parse it
	parser := tasks.ParseTask("", "", time.Now())
	in := parse.NewInput(fmt.Sprintf("- [%s] %s", m.status.Value(), m.text.Value()))
	task, ok, _ := parser.Parse(in)

	// TODO: Better parsing is probably needed for the IsValid field to be useful.
	if !ok {
		m.text.IsValid = false
		return
	}
	m.text.IsValid = true

	// If it parses, update the fields
	m.fields.Due = task.Due()
	m.fields.Scheduled = task.Scheduled()
	m.fields.Priority = task.Priority()
	m.fields.Every = task.Every()
	m.fields.Completed = m.computeCompleted(task)
}

func (m *Model) computeCompleted(task tasks.Task) *time.Time {
	completed := func(t time.Time) string {
		return fmt.Sprintf(" completed:%s", t.Format("2006-01-02"))
	}

	// If fields completed is set, but one of text or status is not then we are removing the completed value.
	// Ensure both text and status are both set to not completed.
	if m.fields.Completed != nil {
		// If all three are set then nothing has changed.
		if strings.Contains(m.text.Value(), completed(*m.fields.Completed)[1:]) && m.status.Value() == tasks.Done {
			return m.fields.Completed
		}
		if strings.Contains(m.text.Value(), completed(*m.fields.Completed)) {
			cursor := m.text.Cursor()
			m.text.SetValue(strings.ReplaceAll(m.text.Value(), completed(*m.fields.Completed), ""))
			m.text.SetCursor(cursor)
		}
		m.status.SetValue(tasks.Todo)
		return nil
	}

	// If one of text or status is set but the fields completed is not, then we are adding the completed value.
	// Ensure both text and status are both set to completed.
	if m.fields.Completed == nil {
		// If all three are unset then nothing has changed.
		if m.status.Value() != tasks.Done && !strings.Contains(m.text.Value(), "completed:") {
			return nil
		}
		m.status.SetValue(tasks.Done)

		// If the text does not contain the completed value, set it to the current time.
		if !strings.Contains(m.text.Value(), "completed:") {
			now := time.Now()
			cursor := m.text.Cursor()
			m.text.SetValue(m.text.Value() + completed(now))
			m.text.SetCursor(cursor)
			return &now
		}

		// Otherwise just return the parsed value (the one set by the text).
		return task.Completed()

	}

	// We should never reach this point. Not sure why go can't work this out itself...
	return m.fields.Completed
}
