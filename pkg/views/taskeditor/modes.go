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
	"github.com/notedownorg/notedown/pkg/ast"
	"github.com/notedownorg/task/pkg/components/statusbar"
)

type Mode func(*Model)

func WithAdd(status ast.Status, text string) Mode {
	return func(m *Model) {
		m.status = NewStatus(m.ctx, status).Focus()
		m.text = NewText(m.ctx).SetValue(text)
		m.footer = statusbar.New(m.ctx, statusbar.NewMode("add task", statusbar.ActionCreate), m.tasks)
		m.fields = NewFields(m.ctx)
		m.text.SetCursor(0)
		m.parseTask()
	}
}

func WithEdit(task ast.Task) Mode {
	return func(m *Model) {
		m.status = NewStatus(m.ctx, task.Status()).Focus()
		m.text = NewText(m.ctx).SetValue(task.Body())
		m.footer = statusbar.New(m.ctx, statusbar.NewMode("edit task", statusbar.ActionEdit), m.tasks)
		m.fields = NewFields(m.ctx)
		m.text.SetCursor(0)
		m.parseTask()
	}
}
