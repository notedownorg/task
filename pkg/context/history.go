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

import tea "github.com/charmbracelet/bubbletea/v2"

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
