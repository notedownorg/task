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

package projectmanager

import "github.com/charmbracelet/bubbles/v2/key"

type KeyMap struct {
	ToggleFocus key.Binding

	CursorUp   key.Binding
	CursorDown key.Binding

	AddTask        key.Binding
	EditTask       key.Binding
	DeleteTask     key.Binding
	RescheduleTask key.Binding
	CompleteTask   key.Binding
}

var DefaultKeyMap = KeyMap{
	ToggleFocus: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "toggle metadata tasks list and completed tasks list"),
	),
	CursorUp: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "move cursor up"),
	),
	CursorDown: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move cursor down"),
	),
	AddTask: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add a new task"),
	),
	EditTask: key.NewBinding(
		key.WithKeys("e", "enter"),
		key.WithHelp("e/enter", "edit the selected task"),
	),
	DeleteTask: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete the selected task"),
	),
	RescheduleTask: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reschedule the selected task"),
	),
	CompleteTask: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "complete the selected task"),
	),
}
