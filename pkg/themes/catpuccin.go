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

var (
	CatpuccinMocha = Theme{
		Panel: lipgloss.Color("#313244"),

		Text:       lipgloss.Color("#CDD6F4"),
		TextCursor: lipgloss.Color("#11111B"),

		Red:    lipgloss.Color("#F38BA8"),
		Green:  lipgloss.Color("#A6E3A1"),
		Yellow: lipgloss.Color("#F9E2AF"),
		Blue:   lipgloss.Color("#89B4FA"),
	}
)
