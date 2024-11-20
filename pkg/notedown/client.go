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

package notedown

import (
	"time"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/fileserver/writer"
	"github.com/notedownorg/notedown/pkg/providers/daily"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
)

type TaskReader interface {
	ListTasks(tasks.Fetcher, ...tasks.ListOption) []tasks.Task
	TaskSummary() int
}

type TaskWriter interface {
	CreateTask(
		string,
		int,
		string,
		tasks.Status,
		...tasks.TaskOption,
	) error
	UpdateTask(tasks.Task) error
	DeleteTask(tasks.Task) error
}

type DailyReader interface {
}

type DailyWriter interface {
	EnsureDaily(time.Time, time.Duration) (daily.Daily, bool, error)
}

type Client interface {
	TaskReader
	TaskWriter
	DailyReader
	DailyWriter
	Subscribe(chan tasks.Event)
}

type client struct {
	*tasks.TaskClient
	*daily.DailyClient
}

func NewClient(root string) (Client, error) {
	read, err := reader.NewClient(root, "task")
	if err != nil {
		return nil, err
	}
	write := writer.NewClient(root)

	// Configure the task client
	taskReaderChannel := make(chan reader.Event)
	read.Subscribe(taskReaderChannel, reader.WithInitialDocuments())
	tasksClient := tasks.NewClient(write, taskReaderChannel, tasks.WithInitialLoadWaiter(100*time.Millisecond))

	// Configure the daily note client
	dailyReaderChannel := make(chan reader.Event)
	read.Subscribe(dailyReaderChannel, reader.WithInitialDocuments())
	dailyClient := daily.NewClient(write, dailyReaderChannel, daily.WithInitialLoadWaiter(100*time.Millisecond))

	return &client{
		TaskClient:  tasksClient,
		DailyClient: dailyClient,
	}, nil
}

func (c *client) Subscribe(sub chan tasks.Event) {
	c.TaskClient.Subscribe(sub)
}
