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
	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/notedownorg/notedown/pkg/providers/projects"
	"github.com/notedownorg/task/pkg/context"
)

var _ context.Listener = &ProjectListener{}

type ProjectListener struct {
	ch <-chan projects.Event
}

type ProjectEvent struct{}

func NewProjectListener(ch <-chan projects.Event) *ProjectListener {
	return &ProjectListener{ch: ch}
}

func (l *ProjectListener) Init() tea.Cmd {
	return func() tea.Msg {
		return ProjectEvent{}
	}
}

func (l *ProjectListener) Receive(msg tea.Msg) tea.Cmd {
	msg, ok := msg.(ProjectEvent)

	// If it's not a ProjectEvent, we don't care about it
	if !ok {
		return nil
	}

	// If it is a ProjectEvent, we need to create a new command that waits for the next project event to come in
	// and then responds with a new ProjectEvent message. This in turn will trigger a UI refresh/re-render when it resolves.
	return func() tea.Msg {
		slog.Debug("project listener waiting for next project event")
		<-l.ch
		slog.Debug("project listener received project event, sending project refresh message")
		return ProjectEvent{}
	}
}
