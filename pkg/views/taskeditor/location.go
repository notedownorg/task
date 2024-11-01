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

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/model"
)

type Location struct {
	base model.Base

	ctx *context.ProgramContext

	file string
	line int
}

func NewLocation(ctx *context.ProgramContext) *Location {
	return &Location{
		ctx: ctx,
	}
}

func (l *Location) Init() (tea.Model, tea.Cmd) {
	return l, nil
}

func (l *Location) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return l, nil
}

func (l *Location) SetLocation(file string, line int) *Location {
	l.file = file
	l.line = line
	return l
}

func (l *Location) View() string {
	file := l.file
	line := fmt.Sprintf("%d", l.line)
	if file == "" {
		file = "<select>"
		line = "At End"
	}
	return l.base.NewStyle().
		Render(fmt.Sprintf("  %s 󰁕 %s", file, line))
}
