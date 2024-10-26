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

package listeners

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/notedownorg/notedown/pkg/workspace/tasks"
	"github.com/notedownorg/task/pkg/context"
)

var _ context.Listener = &TaskListener{}

type TaskListener struct {
	ch <-chan tasks.Event
}

type TaskEvent struct{}

func NewTaskListener(ch <-chan tasks.Event) *TaskListener {
	return &TaskListener{ch: ch}
}

func (l *TaskListener) Init() tea.Cmd {
	return func() tea.Msg {
		return TaskEvent{}
	}
}

func (l *TaskListener) Receive(msg tea.Msg) tea.Cmd {
	msg, ok := msg.(TaskEvent)

	// If it's not a TaskEvent, we don't care about it
	if !ok {
		return nil
	}

	// If it is a TaskEvent, we need to create a new command that waits for the next task event to come in
	// and then responds with a new TaskEvent message. This in turn will trigger a UI refresh/re-render when it resolves.
	return func() tea.Msg {
		<-l.ch
		return TaskEvent{}
	}
}
