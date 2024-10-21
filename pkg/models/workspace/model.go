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

package workspace

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/notedownorg/notedown/pkg/workspace/tasks"
)

var _ tea.Model = &Model{}

func New(tasks *tasks.Client) *Model {
	return &Model{
		tasks: tasks,
	}
}

type Model struct {
	tasks *tasks.Client
}

func (m *Model) Init() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *Model) View() string {
	allTasks, _ := m.tasks.ListTasks(tasks.FetchAllTasks())
	return fmt.Sprintf("%d tasks loaded", len(allTasks))
}
