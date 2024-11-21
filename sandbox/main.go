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
	"os"

	"github.com/spf13/pflag"
)

var (
	fileCount = pflag.Int("files", 30, "number of files to generate (2/3rds projects, 1/3rd empty files)")
	taskCount = pflag.Int("tasks", 300, "number of tasks to generate")
)

func main() {
	pflag.Parse()
	tmp, _ := os.MkdirTemp("", "task-sandbox")

	fmt.Println(tmp)
	GenerateWorkspace(tmp, *fileCount, *taskCount)
}
