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

package model

import "github.com/charmbracelet/lipgloss"

// Base provides basic helper functionality for all (non-view) models.
// It should be a field in all models but cannot be embedded in the idiomatic Go manner.
// This has the advantage of the parent being able to control which of the Base methods are exposed.
// Maybe it is possible to do embedded though? Open to suggestions!
type Base struct {
	width int
	h     int

	margin []int
}

// Returns a new lipgloss.Style with all base fields set.
func (b Base) NewStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(b.width).
		Height(b.h).
		Margin(b.margin...)
}

// Pass through to lipgloss.Style Height
// Height sets the height of the block before applying margins. If the height of the text block is less than this value after applying padding (or not), the block will be set to this height.
func (b *Base) Height(height int) *Base {
	b.h = height
	return b
}

// Pass through to lipgloss.Style Width
// Width sets the width of the block before applying margins. The width, if set, also determines where text will wrap.
func (b *Base) Width(width int) *Base {
	b.width = width
	return b
}

// AvailableWidth returns the width of the block minus the margins.
func (b Base) AvailableWidth() int {
	res := b.width
	if len(b.margin) == 1 {
		return res - b.margin[0]*2
	}
	if len(b.margin) == 2 || len(b.margin) == 3 {
		return res - b.margin[1]*2
	}
	if len(b.margin) == 4 {
		return res - b.margin[1] + b.margin[3]
	}
	return res
}

// Pass through to lipgloss.Style Margin
// Margin is a shorthand method for setting margins on all sides at once.
//
// With one argument, the value is applied to all sides.
//
// With two arguments, the value is applied to the vertical and horizontal sides, in that order.
//
// With three arguments, the value is applied to the top side, the horizontal sides, and the bottom side, in that order.
//
// With four arguments, the value is applied clockwise starting from the top side, followed by the right side, then the bottom, and finally the left.
//
// With more than four arguments no margin will be added.
func (b *Base) Margin(i ...int) *Base {
	b.margin = i
	return b
}
