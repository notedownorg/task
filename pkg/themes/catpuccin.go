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
		Panel:       CatpuccinMochaPalette.Surface0,
		BorderFaint: CatpuccinMochaPalette.Surface1,

		Text:       CatpuccinMochaPalette.Text,
		TextCursor: CatpuccinMochaPalette.Crust,

		Red:    CatpuccinMochaPalette.Red,
		Green:  CatpuccinMochaPalette.Green,
		Yellow: CatpuccinMochaPalette.Yellow,
		Blue:   CatpuccinMochaPalette.Blue,
	}

	// https://catppuccin.com/palette
	CatpuccinMochaPalette = CatpuccinPalette{
		Rosewater: lipgloss.Color("#F5E0DC"),
		Flamingo:  lipgloss.Color("#F2CDCD"),
		Pink:      lipgloss.Color("#F5C2E7"),
		Mauve:     lipgloss.Color("#CBA6F7"),
		Red:       lipgloss.Color("#F38BA8"),
		Maroon:    lipgloss.Color("#EBA0AC"),
		Peach:     lipgloss.Color("#FAB387"),
		Yellow:    lipgloss.Color("#F9E2AF"),
		Green:     lipgloss.Color("#A6E3A1"),
		Teal:      lipgloss.Color("#94E2D5"),
		Sky:       lipgloss.Color("#89DCEB"),
		Sapphire:  lipgloss.Color("#74C7EC"),
		Blue:      lipgloss.Color("#89B4FA"),
		Lavender:  lipgloss.Color("#B4BEFE"),

		Text:     lipgloss.Color("#CDD6F4"),
		Subtext1: lipgloss.Color("#BAC2DE"),
		Subtext0: lipgloss.Color("#A6ADC8"),

		Overlay2: lipgloss.Color("#9399B2"),
		Overlay1: lipgloss.Color("#7F849C"),
		Overlay0: lipgloss.Color("#6C7086"),

		Surface2: lipgloss.Color("#585B70"),
		Surface1: lipgloss.Color("#45475A"),
		Surface0: lipgloss.Color("#313244"),

		Base:   lipgloss.Color("#1E1E2E"),
		Mantle: lipgloss.Color("#181825"),
		Crust:  lipgloss.Color("#11111B"),
	}
)

type CatpuccinPalette struct {
	Rosewater lipgloss.Color
	Flamingo  lipgloss.Color
	Pink      lipgloss.Color
	Mauve     lipgloss.Color
	Red       lipgloss.Color
	Maroon    lipgloss.Color
	Peach     lipgloss.Color
	Yellow    lipgloss.Color
	Green     lipgloss.Color
	Teal      lipgloss.Color
	Sky       lipgloss.Color
	Sapphire  lipgloss.Color
	Blue      lipgloss.Color
	Lavender  lipgloss.Color

	Text     lipgloss.Color
	Subtext1 lipgloss.Color
	Subtext0 lipgloss.Color

	Overlay2 lipgloss.Color
	Overlay1 lipgloss.Color
	Overlay0 lipgloss.Color

	Surface2 lipgloss.Color
	Surface1 lipgloss.Color
	Surface0 lipgloss.Color

	Base   lipgloss.Color
	Mantle lipgloss.Color
	Crust  lipgloss.Color
}
