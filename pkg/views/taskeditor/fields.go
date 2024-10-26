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
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/notedownorg/notedown/pkg/ast"
	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/model"
)

type Fields struct {
	base model.Base

	ctx *context.ProgramContext

	due       *time.Time
	scheduled *time.Time
	completed *time.Time
	priority  *int
	every     *ast.Every
}

func NewFields(ctx *context.ProgramContext) *Fields {
	return &Fields{
		ctx: ctx,
	}
}

func (f *Fields) Init() (tea.Model, tea.Cmd) {
	return f, nil
}

func (f *Fields) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return f, nil
}

func (f Fields) unset() bool {
	return f.due == nil && f.scheduled == nil && f.completed == nil && f.priority == nil && f.every == nil
}

func (s *Fields) View() string {
	if s.unset() {
		return lipgloss.NewStyle().
			Padding(0, 1).
			Background(s.ctx.Theme.Panel).
			Render("no fields set")
	}

	return s.base.NewStyle().
		Render("todo")
}
