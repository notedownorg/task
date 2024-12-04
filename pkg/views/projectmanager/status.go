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

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/notedownorg/notedown/pkg/providers/projects"
	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/model"
	"github.com/notedownorg/task/pkg/styling/colors"
	"github.com/notedownorg/task/pkg/styling/icons"
)

type Status struct {
	base model.Base

	ctx *context.ProgramContext

	value   projects.Status
	focused bool
}

func NewStatus(ctx *context.ProgramContext, value projects.Status) *Status {
	return &Status{
		ctx:   ctx,
		value: value,
	}
}

func (s *Status) Init() (tea.Model, tea.Cmd) {
	return s, nil
}

func (s *Status) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !s.focused {
		return s, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ", "space":
			s.value = projects.Backlog
		case "b", "B":
			s.value = projects.Blocked
		case "/":
			s.value = projects.Active
		case "x", "X":
			s.value = projects.Archived
		case "a", "A":
			s.value = projects.Abandoned
		}
	}

	return s, nil
}

func (s *Status) Focus() *Status {
	s.focused = true
	return s
}

func (s *Status) Blur() *Status {
	s.focused = false
	return s
}

func (s Status) Value() projects.Status {
	return s.value
}

func (s *Status) SetValue(value projects.Status) {
	s.value = value
}

func (s *Status) Margin(i ...int) *Status {
	s.base.Margin(i...)
	return s
}

func (s *Status) View() string {
	_, c, _ := colors.Project(s.ctx.Theme, s.value)
	return s.base.NewStyle().
		Foreground(c).
		Blink(s.focused).
		Render(icons.Project(s.value))
}
