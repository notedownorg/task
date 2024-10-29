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

package themes

import (
	"github.com/charmbracelet/lipgloss"
)

// TODO: The colors should probably have more generic/correct(?) names
type Theme struct {
	Panel       lipgloss.Color
	BorderFaint lipgloss.Color

	Text       lipgloss.Color
	TextCursor lipgloss.Color
	TextFaint  lipgloss.Color

	// Basic Terminal Colors
	Red     lipgloss.Color
	Green   lipgloss.Color
	Yellow  lipgloss.Color
	Blue    lipgloss.Color
	Magenta lipgloss.Color
	Cyan    lipgloss.Color

	// Additional Colors
	Orange   lipgloss.Color
	Pink     lipgloss.Color
	RedSoft  lipgloss.Color
	BlueSoft lipgloss.Color
}
