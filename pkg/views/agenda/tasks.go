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

package agenda

import (
	"github.com/notedownorg/notedown/pkg/ast"
	"github.com/notedownorg/notedown/pkg/workspace/tasks"
	"github.com/notedownorg/task/pkg/components/groupedlist"
)

func (m *Model) visibleTasks() []groupedlist.Group[ast.Task] {
	prev := m.date.Add(-1)
	next := m.date.AddDate(0, 0, 1).Add(-1)

	overdue := m.tasks.ListTasks(
		tasks.FetchAllTasks(),
		tasks.WithFilters(
			tasks.FilterByStatus(ast.Todo, ast.Doing, ast.Blocked),
			tasks.FilterByDueDate(nil, &next),
		),
		tasks.WithSorters(
			tasks.SortByStatus(tasks.AgendaOrder()),
			tasks.SortByPriority(),
		),
	)

	doing := groupedlist.Group[ast.Task]{Name: statusName[ast.Doing], Items: make([]ast.Task, 0)}
	todo := groupedlist.Group[ast.Task]{Name: statusName[ast.Todo], Items: make([]ast.Task, 0)}
	blocked := groupedlist.Group[ast.Task]{Name: statusName[ast.Blocked], Items: make([]ast.Task, 0)}

	for _, t := range overdue {
		switch t.Status() {
		case ast.Doing:
			doing.Items = append(doing.Items, t)
		case ast.Todo:
			todo.Items = append(todo.Items, t)
		case ast.Blocked:
			blocked.Items = append(blocked.Items, t)
		}
	}

	done := m.tasks.ListTasks(
		tasks.FetchAllTasks(),
		tasks.WithFilters(
			tasks.FilterByStatus(ast.Done),
			tasks.FilterByCompletedDate(&prev, &next),
		),
		tasks.WithSorters(tasks.SortByAlphabetical()),
	)

	return []groupedlist.Group[ast.Task]{doing, todo, blocked, {Name: statusName[ast.Done], Items: done}}
}

var statusName = map[ast.Status]string{
	ast.Todo:      "Todo",
	ast.Doing:     "Doing",
	ast.Blocked:   "Blocked",
	ast.Done:      "Done",
	ast.Abandoned: "Abandoned",
}
