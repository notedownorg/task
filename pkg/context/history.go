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
	"log/slog"

	tea "github.com/charmbracelet/bubbletea/v2"
)

// History is a stack of models that the user has navigated through.
type History struct {
	Items []tea.Model
}

func (h *History) Push(m tea.Model) {
	h.Items = append(h.Items, m)
}

func (h *History) Pop() (tea.Model, bool) {
	if len(h.Items) == 0 {
		return nil, false
	}

	m := h.Items[len(h.Items)-1]
	h.Items = h.Items[:len(h.Items)-1]
	return m, true
}

func (h History) Peek() (tea.Model, bool) {
	if len(h.Items) == 0 {
		return nil, false
	}
	return h.Items[len(h.Items)-1], true
}

func (h History) Len() int {
	return len(h.Items)
}

type NavigationEvent struct{}

func (c *ProgramContext) Back() tea.Model {
	// As history includes the current model, check that we are not at the beginning of the history.
	// Then pop the current and return the previous
	if c.History.Len() > 1 {
		c.History.Pop()
		m, _ := c.History.Peek()

		// Call update before returning the model to allow the view to update itself.
		m.Update(&NavigationEvent{})
		return m
	}

	// Return nil if there is no history to go back to.
	slog.Debug("we've reached the beginning of the navigation history so there is no view to navigate back to")
	return nil
}

func (c *ProgramContext) Navigate(next tea.Model) (tea.Model, tea.Cmd) {
	c.History.Push(next)
	return next, nil
}
