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

package groupedlist

import (
	"github.com/charmbracelet/bubbles/v2/viewport"
	"github.com/charmbracelet/lipgloss"
)

type Group[T any] struct {
	Name  string
	Items []T
}

// For each group in the list the following occurs:
//  1. Header is rendered
//  2. Each item is rendered with item or selected depending on if it is the cursor
//  3. Footer is rendered
//
// If you don't want a header or footer, leave the function nil.
// If there are no items in the group it is skipped entirely.
type Renderers[T any] struct {
	Header   func(string, int) string
	Footer   func(string, int) string
	Item     func(T, int) string
	Selected func(T, int) string
}

type Model[T any] struct {
	groups     []Group[T]
	totalItems int
	renderers  Renderers[T]

	focus          bool
	cursor         int // index of the selected item
	cursorAbsolute int // index of the position in the viewport content
	viewport       viewport.Model
}

func New[T any](opts ...Option[T]) *Model[T] {
	m := &Model[T]{
		groups:   make([]Group[T], 0),
		viewport: viewport.New(0, 20),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Model[T]) SetGroups(groups []Group[T]) {
	m.groups = groups
	m.totalItems = 0
	for _, g := range groups {
		m.totalItems += len(g.Items)
	}

	// Reset the cursor if it's now out of bounds
	m.cursor = clamp(m.cursor, 0, m.totalItems-1)

	m.updateViewport()

	// Ensure the viewport finishes at the bottom of the content (no trailing whitespace)
	if m.viewport.YOffset > m.viewport.TotalLineCount()-m.viewport.Height {
		m.viewport.SetYOffset(m.viewport.TotalLineCount() - m.viewport.Height)
	}

	// If the cursor is now outside of the visible viewport, adjust the YOffset to make it visible
	if m.cursorAbsolute > m.viewport.YOffset+m.viewport.Height-1 || m.cursorAbsolute < m.viewport.YOffset {
		m.viewport.SetYOffset(m.cursorAbsolute - m.viewport.Height/2)
	}
}

func (m *Model[T]) MoveUp(n int) {
	m.cursor = clamp(m.cursor-n, 0, m.totalItems-1)
	m.updateViewport()

	// TODO: Come up with a better way to handle scrolling?
	// Start scrolling if the cursor is roughly above the middle of the viewport
	if m.cursorAbsolute < m.viewport.TotalLineCount()-(m.viewport.Height/2) {
		m.viewport.SetYOffset(m.viewport.YOffset - n)
	}
}

func (m *Model[T]) MoveDown(n int) {
	m.cursor = clamp(m.cursor+n, 0, m.totalItems-1)
	m.updateViewport()

	// TODO: Come up with a better way to handle scrolling?
	// Start scrolling if the cursor is roughly below the middle of the viewport
	if m.cursorAbsolute > m.viewport.Height/2 {
		m.viewport.SetYOffset(m.viewport.YOffset + n)
	}
}

func (m *Model[T]) Focus() *Model[T] {
	m.focus = true
	m.updateViewport()
	return m
}

func (m Model[T]) Focused() bool {
	return m.focus
}

func (m *Model[T]) Blur() *Model[T] {
	m.focus = false
	m.updateViewport()
	return m
}

func (m *Model[T]) Width(i int) *Model[T] {
	m.viewport.Width = i
	m.updateViewport()
	return m
}

func (m *Model[T]) Height(i int) *Model[T] {
	m.viewport.Height = i
	m.updateViewport()
	return m
}

// there are various optimizations that could be made here
// right now we just re-render everything
func (m *Model[T]) updateViewport() {
	renderedLines := make([]string, 0)

	groupIndex := 0
	itemIndex := 0
	for {
		if groupIndex >= len(m.groups) {
			break
		}
		group := m.groups[groupIndex]
		if len(group.Items) != 0 {
			renderedLines = m.renderHeader(renderedLines, groupIndex)
			for i := 0; i < len(group.Items); i++ {
				if itemIndex == m.cursor && m.focus {
					renderedLines = m.renderSelected(renderedLines, groupIndex, i)
					m.cursorAbsolute = len(renderedLines) - 1
				} else {
					renderedLines = m.renderItem(renderedLines, groupIndex, i)
				}
				itemIndex++
			}
			renderedLines = m.renderFooter(renderedLines, groupIndex)
		}
		groupIndex++
	}

	m.viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, renderedLines...),
	)
}

func (m Model[T]) renderHeader(acc []string, group int) []string {
	if m.renderers.Header == nil {
		return acc
	}
	return append(acc, m.renderers.Header(m.groups[group].Name, m.viewport.Width))
}

func (m Model[T]) renderFooter(acc []string, group int) []string {
	if m.renderers.Footer == nil {
		return acc
	}
	return append(acc, m.renderers.Footer(m.groups[group].Name, m.viewport.Width))
}

func (m Model[T]) renderItem(acc []string, group int, index int) []string {
	return append(acc, m.renderers.Item(m.groups[group].Items[index], m.viewport.Width))
}

func (m Model[T]) renderSelected(acc []string, group int, index int) []string {
	return append(acc, m.renderers.Selected(m.groups[group].Items[index], m.viewport.Width))
}

func (m Model[T]) View() string {
	return m.viewport.View()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}
