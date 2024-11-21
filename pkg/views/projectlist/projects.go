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

package projectlist

import (
	"github.com/notedownorg/notedown/pkg/providers/projects"
	"github.com/notedownorg/task/pkg/components/groupedlist"
)

func (m *Model) updateProjects() {
	open := m.nd.ListProjects(
		projects.FetchAllProjects(),
		projects.WithFilter(
			projects.FilterByStatus(projects.Active, projects.Backlog, projects.Blocked, projects.Archived),
		),
		projects.WithSorters(), // empty defaults to alphabetical
	)

	active := groupedlist.Group[projects.Project]{
		Name:  statusName[projects.Active],
		Items: projects.WithFilter(projects.FilterByStatus(projects.Active))(open),
	}
	backlog := groupedlist.Group[projects.Project]{
		Name:  statusName[projects.Backlog],
		Items: projects.WithFilter(projects.FilterByStatus(projects.Backlog))(open),
	}
	blocked := groupedlist.Group[projects.Project]{
		Name:  statusName[projects.Blocked],
		Items: projects.WithFilter(projects.FilterByStatus(projects.Blocked))(open),
	}
	archived := groupedlist.Group[projects.Project]{
		Name:  statusName[projects.Archived],
		Items: projects.WithFilter(projects.FilterByStatus(projects.Archived))(open),
	}

	m.projectlist.SetGroups([]groupedlist.Group[projects.Project]{active, backlog, blocked})
	m.closed.SetGroups([]groupedlist.Group[projects.Project]{archived})

}

var statusName = map[projects.Status]string{
	projects.Active:   "Active",
	projects.Backlog:  "Backlog",
	projects.Blocked:  "Blocked",
	projects.Archived: "Archived",
}
