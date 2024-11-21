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
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/notedownorg/task/pkg/components/groupedlist"
)

func (m *Model) updateTasks() {
	prev := m.date.Add(-1)
	next := m.date.AddDate(0, 0, 1).Add(-1)

	overdue := m.nd.ListTasks(
		tasks.FetchAllTasks(),
		tasks.WithFilter(
			tasks.And(
				tasks.FilterByStatus(tasks.Todo, tasks.Doing, tasks.Blocked),
				tasks.Or(
					tasks.FilterByDueDate(nil, &next),
					tasks.FilterByScheduledDate(nil, &next),
				),
			),
		),
		tasks.WithSorters(
			tasks.SortByStatus(tasks.AgendaOrder()),
			tasks.SortByPriority(),
		),
	)

	doing := groupedlist.Group[tasks.Task]{
		Name:  statusName[tasks.Doing],
		Items: tasks.WithFilter(tasks.FilterByStatus(tasks.Doing))(overdue),
	}
	todo := groupedlist.Group[tasks.Task]{
		Name:  statusName[tasks.Todo],
		Items: tasks.WithFilter(tasks.FilterByStatus(tasks.Todo))(overdue),
	}
	blocked := groupedlist.Group[tasks.Task]{
		Name:  statusName[tasks.Blocked],
		Items: tasks.WithFilter(tasks.FilterByStatus(tasks.Blocked))(overdue),
	}

	m.tasklist.SetGroups([]groupedlist.Group[tasks.Task]{doing, todo, blocked})

	done := m.nd.ListTasks(
		tasks.FetchAllTasks(),
		tasks.WithFilter(
			tasks.And(
				tasks.FilterByStatus(tasks.Done),
				tasks.FilterByCompletedDate(&prev, &next),
			),
		),
		tasks.WithSorters(), // empty defaults to alphabetical
	)
	m.completed.SetGroups([]groupedlist.Group[tasks.Task]{{Name: "Completed", Items: done}})
}

var statusName = map[tasks.Status]string{
	tasks.Todo:      "Todo",
	tasks.Doing:     "Doing",
	tasks.Blocked:   "Blocked",
	tasks.Done:      "Done",
	tasks.Abandoned: "Abandoned",
}
