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

package taskreschedule

import "github.com/charmbracelet/bubbles/v2/key"

type KeyMap struct {
	Today       key.Binding
	Tommorrow   key.Binding
	InTwoDays   key.Binding
	InThreeDays key.Binding
	InFourDays  key.Binding
	InFiveDays  key.Binding
	InSixDays   key.Binding
	InSevenDays key.Binding
	InFortnight key.Binding
	NextMonth   key.Binding
	NextYear    key.Binding
}

var DefaultKeyMap = KeyMap{
	Today: key.NewBinding(
		key.WithKeys("0"),
		key.WithHelp("0", "reschedule to today"),
	),
	Tommorrow: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "reschedule to tomorrow"),
	),
	InTwoDays: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "reschedule in two days"),
	),
	InThreeDays: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "reschedule in three days"),
	),
	InFourDays: key.NewBinding(
		key.WithKeys("4"),
		key.WithHelp("4", "reschedule in four days"),
	),
	InFiveDays: key.NewBinding(
		key.WithKeys("5"),
		key.WithHelp("5", "reschedule in five days"),
	),
	InSixDays: key.NewBinding(
		key.WithKeys("6"),
		key.WithHelp("6", "reschedule in six days"),
	),
	InSevenDays: key.NewBinding(
		key.WithKeys("7"),
		key.WithHelp("7", "reschedule in seven days"),
	),
	InFortnight: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "reschedule in a fortnight (fourteen days)"),
	),
	NextMonth: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "reschedule to 1st of next month"),
	),
	NextYear: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "reschedule to 1st of next year"),
	),
}
