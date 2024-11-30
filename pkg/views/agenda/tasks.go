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
	"time"

	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/notedownorg/task/pkg/components/groupedlist"
	"github.com/notedownorg/task/pkg/notedown"
)

func (m *Model) updateTasks() {
	due := due(m.nd, m.date)
	done := done(m.nd, m.date)

	doing := groupedlist.Group[tasks.Task]{
		Name:  statusName[tasks.Doing],
		Items: tasks.WithFilter(tasks.FilterByStatus(tasks.Doing))(due),
	}
	todo := groupedlist.Group[tasks.Task]{
		Name:  statusName[tasks.Todo],
		Items: tasks.WithFilter(tasks.FilterByStatus(tasks.Todo))(due),
	}
	blocked := groupedlist.Group[tasks.Task]{
		Name:  statusName[tasks.Blocked],
		Items: tasks.WithFilter(tasks.FilterByStatus(tasks.Blocked))(due),
	}

	m.tasklist.SetGroups([]groupedlist.Group[tasks.Task]{doing, todo, blocked})
	m.completed.SetGroups([]groupedlist.Group[tasks.Task]{{Name: "Completed", Items: done}})
}

func due(nd notedown.Client, date time.Time) []tasks.Task {
	next := time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, 0, date.Location()).Add(-time.Second)

	return nd.ListTasks(
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
}

func done(nd notedown.Client, date time.Time) []tasks.Task {
	prev := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	next := time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, 0, date.Location()).Add(-time.Second)

	return nd.ListTasks(
		tasks.FetchAllTasks(),
		tasks.WithFilter(
			tasks.And(
				tasks.FilterByStatus(tasks.Done),
				tasks.FilterByCompletedDate(&prev, &next),
			),
		),
		tasks.WithSorters(), // empty defaults to alphabetical
	)
}

var statusName = map[tasks.Status]string{
	tasks.Todo:      "Todo",
	tasks.Doing:     "Doing",
	tasks.Blocked:   "Blocked",
	tasks.Done:      "Done",
	tasks.Abandoned: "Abandoned",
}
