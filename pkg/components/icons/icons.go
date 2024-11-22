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

package icons

import (
	"log/slog"

	"github.com/notedownorg/notedown/pkg/providers/tasks"
)

func Task(status tasks.Status) string {
	switch status {
	case tasks.Todo:
		return ""
	case tasks.Blocked:
		return ""
	case tasks.Doing:
		return ""
	case tasks.Done:
		return ""
	case tasks.Abandoned:
		return ""
	default:
		slog.Warn("unknown task status", "status", status)
		return " "
	}
}
