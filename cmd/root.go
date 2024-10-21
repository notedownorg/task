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

package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/notedownorg/notedown/pkg/workspace/documents/reader"
	"github.com/notedownorg/notedown/pkg/workspace/documents/writer"
	"github.com/notedownorg/notedown/pkg/workspace/tasks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/notedownorg/task/pkg/models/workspace"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "A task management CLI & TUI for your Notedown notes",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := loadConfig()

		reader, err := reader.NewClient(cfg.root, "task")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		writer := writer.NewClient(cfg.root)

		tasksClient := tasks.NewClient(writer, reader.Subscribe())

		ws := workspace.New(tasksClient)

		p := tea.NewProgram(ws, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

type config struct {
	root string
}

func loadConfig() config {
	cfg := config{}
	cfg.root = viper.GetString("dir")
	if cfg.root == "" {
		fmt.Println("Please set NOTEDOWN_DIR environment variable to the root of your Notedown workspace")
		os.Exit(1)
	}
	return cfg
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("notedown")
	viper.BindEnv("dir")
	viper.AutomaticEnv() // read in environment variables that match
}
