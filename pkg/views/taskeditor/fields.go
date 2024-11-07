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
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/model"
)

type Fields struct {
	base model.Base

	ctx *context.ProgramContext

	Due       *time.Time
	Scheduled *time.Time
	Completed *time.Time
	Priority  *int
	Every     *tasks.Every
	Name      string
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
	return f.Due == nil && f.Scheduled == nil && f.Completed == nil && f.Priority == nil && f.Every == nil
}

func (f Fields) pillStyle(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(color).
		Foreground(f.ctx.Theme.TextCursor).
		Padding(0, 1)
}

func (f *Fields) View() string {
	if f.unset() {
		return lipgloss.NewStyle().
			Padding(0, 1).
			Background(f.ctx.Theme.Panel).
			Render("no fields set")
	}

	var fields []string
	if f.Due != nil {
		fields = append(fields, f.pillStyle(f.ctx.Theme.Green).Render("󰃭 "+f.Due.Format("2006-01-02")))
	}
	if f.Priority != nil {
		fields = append(fields, f.pillStyle(f.ctx.Theme.Yellow).Render("  "+fmt.Sprintf("%d", *f.Priority)))
	}
	if f.Scheduled != nil {
		fields = append(fields, f.pillStyle(f.ctx.Theme.RedSoft).Render("󰀠 "+f.Scheduled.Format("2006-01-02")))
	}
	if f.Every != nil {
		fields = append(fields, f.pillStyle(f.ctx.Theme.Magenta).Render("  "+f.Every.String()))
	}
	if f.Completed != nil {
		fields = append(fields, f.pillStyle(f.ctx.Theme.BlueSoft).Render(" "+f.Completed.Format("2006-01-02")))
	}

	return f.base.NewStyle().
		Render(strings.Join(fields, "  "))
}
