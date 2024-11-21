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

package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/notedownorg/notedown/pkg/fileserver/writer"
	"github.com/notedownorg/notedown/pkg/providers/projects"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/notedownorg/task/pkg/notedown"
	"github.com/tjarratt/babble"
)

func GenerateWorkspace(root string, maxFiles int, maxTasks int) {
	// Create the non-project files before client so we don't have to worry about loading
	files := make([]string, 0, maxFiles)
	babbler := babble.NewBabbler()
	babbler.Count = 2

	for i := 0; i < maxFiles/3; i++ { // 1/3 of the files are empty files
		name := babbler.Babble()
		rel := fmt.Sprintf("%s.md", name)
		files = append(files, rel)
		if err := os.WriteFile(filepath.Join(root, rel), []byte(fmt.Sprintf("# %s\n", name)), 0644); err != nil {
			slog.Error("failed to create file", "file", rel, "error", err)
		}
	}

	// Configure the tasks client
	nd, err := notedown.NewClient(root)
	if err != nil {
		slog.Error("failed to create notedown client", "error", err)
	}

	// Generate projects
	for i := 0; i < (maxFiles/3)*2; i++ { // 2/3rds of the files are projects
		files = append(files, genProject(nd))
	}

	// Create the tasks
	for i := 0; i < maxTasks; i++ {
		genTask(nd, files[rand.Intn(len(files))])
	}

}

var taskStatuses = []tasks.Status{tasks.Todo, tasks.Doing, tasks.Blocked, tasks.Done, tasks.Abandoned}

func genTask(client notedown.Client, file string) {
	opts := []tasks.TaskOption{}

	// Random status
	status := taskStatuses[rand.Intn(len(taskStatuses))]

	// Randomly add other fields

	// Due dates
	switch rand.Intn(4) {
	// no due date for 0
	case 1: // due date in past
		opts = append(opts, tasks.WithDue(time.Now().AddDate(0, 0, -rand.Intn(15))))
	case 2: // due date in future
		opts = append(opts, tasks.WithDue(time.Now().AddDate(0, 0, rand.Intn(15))))
	case 3: // due date today
		opts = append(opts, tasks.WithDue(time.Now()))
	}

	// Scheduled dates
	switch rand.Intn(4) {
	// no scheduled date for 0
	case 1: // scheduled date in past
		opts = append(opts, tasks.WithScheduled(time.Now().AddDate(0, 0, -rand.Intn(15))))
	case 2: // scheduled date in future
		opts = append(opts, tasks.WithScheduled(time.Now().AddDate(0, 0, rand.Intn(15))))
	case 3: // scheduled date today
		opts = append(opts, tasks.WithScheduled(time.Now()))
	}

	// Random priority 0 to 10 or none at all (<0)
	chance := rand.Intn(16) - 6
	if chance > 0 {
		opts = append(opts, tasks.WithPriority(chance))
	}

	// Randomly add everys
	switch rand.Intn(5) - 3 {
	case 0:
		e, _ := tasks.NewEvery("day")
		opts = append(opts, tasks.WithEvery(e))
	case 1:
		e, _ := tasks.NewEvery("week")
		opts = append(opts, tasks.WithEvery(e))

	}

	// If completed we need to set the completed date to a random date in the last 3 days
	if status == tasks.Done {
		opts = append(opts, tasks.WithCompleted(time.Now().AddDate(0, 0, -rand.Intn(3))))
	}

	babbler := babble.NewBabbler()
	babbler.Count = rand.Intn(6) + 1
	babbler.Separator = " "
	if err := client.CreateTask(file, writer.AT_END, fmt.Sprintf("%v", babbler.Babble()), status, opts...); err != nil {
		slog.Error("failed to create task", "file", file, "error", err)
	}
}

var projectStatuses = []projects.Status{projects.Active, projects.Archived, projects.Abandoned, projects.Blocked, projects.Backlog}

func genProject(client notedown.Client) string {
	babbler := babble.NewBabbler()
	babbler.Count = 2
	babbler.Separator = " "
	name := babbler.Babble()
	path := fmt.Sprintf("projects/%s.md", name)

	status := projectStatuses[rand.Intn(len(projectStatuses))]
	if err := client.CreateProject(path, name, status); err != nil {
		slog.Error("failed to create project", "name", name, "error", err)
	}
	return path
}
