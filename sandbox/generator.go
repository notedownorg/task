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

	"github.com/google/uuid"
	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/fileserver/writer"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/tjarratt/babble"
)

func GenerateWorkspace(root string, maxFiles int, maxTasks int) {
	// Create the files before client so we don't have to worry about loading
	files := make([]string, maxFiles)
	for i := 0; i < maxFiles; i++ {
		name := uuid.New().String()
		rel := fmt.Sprintf("%s.md", name)
		files[i] = rel
		if err := os.WriteFile(filepath.Join(root, rel), []byte(fmt.Sprintf("# %s\n", name)), 0644); err != nil {
			slog.Error("failed to create file", "file", rel, "error", err)
		}
	}

	// Configure the tasks client
	read, err := reader.NewClient(root, "task-sandbox-generator")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	write := writer.NewClient(root)
	readSub := make(chan reader.Event)
	read.Subscribe(readSub, reader.WithInitialDocuments())
	tasksClient := tasks.NewClient(write, readSub)

	// Create the tasks
	for i := 0; i < maxTasks; i++ {
		genTask(tasksClient, files[rand.Intn(len(files))])
	}

}

var statuses = []tasks.Status{tasks.Todo, tasks.Doing, tasks.Blocked, tasks.Done, tasks.Abandoned}

func genTask(client *tasks.Client, file string) {
	opts := []tasks.TaskOption{}

	// Random status
	status := statuses[rand.Intn(len(statuses))]

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
	if err := client.Create(file, writer.AT_END, fmt.Sprintf("%v", babbler.Babble()), status, opts...); err != nil {
		slog.Error("failed to create task", "file", file, "error", err)
	}

}
