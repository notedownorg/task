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
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/notedownorg/task/pkg/components/icons"
	"github.com/notedownorg/task/pkg/components/pill"
	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/model"
)

type Fields struct {
	base model.Base

	ctx *context.ProgramContext

	Status    tasks.Status
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

func (f *Fields) View() string {
	theme := f.ctx.Theme

	var fields []string
	fields = append(fields, pill.New(theme.BlueSoft, theme.TextCursor).Render(fmt.Sprintf("%s  %s", icons.Task(f.Status), statusMap[f.Status])))
	if f.Due != nil {
		fields = append(fields, pill.New(theme.Green, theme.TextCursor).Render("󰃭 "+f.Due.Format("2006-01-02")))
	}
	if f.Priority != nil {
		fields = append(fields, pill.New(theme.Yellow, theme.TextCursor).Render("  "+fmt.Sprintf("%d", *f.Priority)))
	}
	if f.Scheduled != nil {
		fields = append(fields, pill.New(theme.RedSoft, theme.TextCursor).Render("󰀠 "+f.Scheduled.Format("2006-01-02")))
	}
	if f.Every != nil {
		fields = append(fields, pill.New(theme.Magenta, theme.TextCursor).Render("󰕇  "+f.Every.String()))
	}
	if f.Completed != nil {
		fields = append(fields, pill.New(theme.Pink, theme.TextCursor).Render(" "+f.Completed.Format("2006-01-02")))
	}

	return f.base.NewStyle().
		Render(strings.Join(fields, "  "))
}

var statusMap = map[tasks.Status]string{
	tasks.Todo:      "todo",
	tasks.Done:      "done",
	tasks.Doing:     "doing",
	tasks.Abandoned: "abandoned",
	tasks.Blocked:   "blocked",
}
