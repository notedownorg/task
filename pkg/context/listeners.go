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

package context

import (
	tea "github.com/charmbracelet/bubbletea/v2"
)

type Listeners []Listener

type Listener interface {
	Init() tea.Cmd
	Receive(msg tea.Msg) tea.Cmd
}

func (l Listeners) Init() tea.Cmd {
	var cmd tea.Cmd
	for _, listener := range l {
		cmd = tea.Batch(cmd, listener.Init())
	}
	return cmd
}

func (l Listeners) Receive(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	for _, listener := range l {
		cmd = tea.Batch(cmd, listener.Receive(msg))
	}
	return cmd
}

func WithListeners(listeners ...Listener) ProgramContextOption {
	return func(p *ProgramContext) {
		p.Listeners = append(p.Listeners, listeners...)
	}
}
