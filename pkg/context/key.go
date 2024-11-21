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

package context

import tea "github.com/charmbracelet/bubbletea/v2"

// GlobalKeyHandlers are used to handle key events at the program level.
// Typically these will be used for program-wide key bindings like quitting or top-level navigation.
type GlobalKeyHandler func(*ProgramContext, tea.KeyMsg) (tea.Model, tea.Cmd)

// SetGlobalKeyHandler replaces the set of global key handlers in the ProgramContext.
// Each handler will be called in order until one returns a non-nil model or command i.e. the handlers are OR'd.
// If you want to do AND logic, you should handle that in a single handler.
func (p *ProgramContext) SetGlobalKeyHandlers(handlers ...GlobalKeyHandler) *ProgramContext {
	p.KeyHandlers = handlers
	return p
}

func HandleQuit() GlobalKeyHandler {
	return func(ctx *ProgramContext, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
		switch msg.String() {
		case "q", "ctrl+c":
			return ctx, tea.Quit
		}
		return nil, nil
	}
}

func HandleBack() GlobalKeyHandler {
	return func(ctx *ProgramContext, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
		switch msg.String() {
		case "esc":
			return ctx.Back(), nil
		}
		return nil, nil
	}
}
