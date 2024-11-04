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
		genTask(tasksClient, files[rand.Intn(len(files))], i)
	}

}

var statuses = []tasks.Status{tasks.Todo, tasks.Doing, tasks.Blocked, tasks.Done, tasks.Abandoned}

func genTask(client *tasks.Client, file string, index int) {
	opts := []tasks.TaskOption{tasks.WithLine(writer.AT_END)}

	// Random status
	status := statuses[rand.Intn(len(statuses))]

	// Randomly add other fields
	// Random due date -1 to +6 days or none at all (-2)
	chance := rand.Intn(9) - 2
	if chance > -2 {
		opts = append(opts, tasks.WithDue(time.Now().AddDate(0, 0, chance)))
	}

	// Random priority 0 to 10 or none at all (-1)
	chance = rand.Intn(11) - 1
	if chance > -1 {
		opts = append(opts, tasks.WithPriority(chance))
	}

	// If completed we need to set the completed date to a random date in the last 3 days
	if status == tasks.Done {
		opts = append(opts, tasks.WithCompleted(time.Now().AddDate(0, 0, -rand.Intn(3))))
	}

	babbler := babble.NewBabbler()
	babbler.Count = rand.Intn(6) + 1
	babbler.Separator = " "
	if err := client.Create(file, fmt.Sprintf("%v", babbler.Babble()), status, opts...); err != nil {
		slog.Error("failed to create task", "file", file, "error", err)
	}

}
