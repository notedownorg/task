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
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/notedownorg/task/pkg/notedown"
)

func TestDue(t *testing.T) {
	// slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug.Level()})))

	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	justBefore := startOfDay.Add(time.Second * -1)
	justAfter := endOfDay.Add(time.Second)
	dates := []time.Time{justBefore, startOfDay, now, endOfDay, justAfter}

	// create a task for each time both due and scheduled
	tsks := make([]tasks.Task, 0)
	for i, due := range dates {
		tsks = append(tsks, tasks.NewTask(
			tasks.NewIdentifier("", "", 0),
			fmt.Sprintf("task %ds", i),
			tasks.Todo,
			tasks.WithDue(due),
		))
		tsks = append(tsks, tasks.NewTask(
			tasks.NewIdentifier("", "", 0),
			fmt.Sprintf("task %dd", i),
			tasks.Todo,
			tasks.WithScheduled(due),
		))
	}
	nd := ndclient(t, tsks...)

	tests := []struct {
		name string
		date time.Time
		want int
	}{
		{
			name: "yesterday",
			date: now.Add(-24 * time.Hour),
			want: 2, // only just before
		},
		{
			name: "today",
			date: now,
			want: 8, // yesterday + start, now, end
		},
		{
			name: "tomorrow",
			date: now.Add(24 * time.Hour),
			want: 10, // all
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := len(due(nd, tt.date))
			if got != tt.want {
				t.Errorf("due() = %v, want %v, all %v", got, tt.want, len(nd.ListTasks(tasks.FetchAllTasks())))
			}
		})
	}
}

func TestDone(t *testing.T) {
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	justBefore := startOfDay.Add(time.Second * -1)
	justAfter := endOfDay.Add(time.Second)
	dates := []time.Time{justBefore, startOfDay, now, endOfDay, justAfter}

	// create a task for each time
	tsks := make([]tasks.Task, 0)
	for i, date := range dates {
		tsks = append(tsks, tasks.NewTask(
			tasks.NewIdentifier("", "", 0),
			fmt.Sprintf("task %d", i),
			tasks.Done,
			tasks.WithCompleted(date),
		))
	}
	nd := ndclient(t, tsks...)

	tests := []struct {
		name string
		date time.Time
		want int
	}{
		{
			name: "yesterday",
			date: now.Add(-24 * time.Hour),
			want: 1,
		},
		{
			name: "today",
			date: now,
			want: 3,
		},
		{
			name: "tomorrow",
			date: now.Add(24 * time.Hour),
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := len(done(nd, tt.date))
			if got != tt.want {
				t.Errorf("done() = %v, want %v, all %v", got, tt.want, len(nd.ListTasks(tasks.FetchAllTasks())))
			}
		})
	}
}

func ndclient(t *testing.T, tsks ...tasks.Task) notedown.Client {
	tmpDir := os.TempDir()

	var b strings.Builder
	b.WriteString("# test\n")
	for _, t := range tsks {
		b.WriteString(fmt.Sprintf("%s\n", t))
	}
	content := b.String()
	if err := os.WriteFile(tmpDir+"/test.md", []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	t.Logf("created %s/test.md", tmpDir)

	nd, err := notedown.NewClient(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	return nd
}
