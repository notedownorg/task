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
	"log/slog"
	"os"
	"path"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/notedownorg/task/pkg/context"
	"github.com/notedownorg/task/pkg/listeners"
	"github.com/notedownorg/task/pkg/notedown"
	"github.com/notedownorg/task/pkg/themes"
	"github.com/notedownorg/task/pkg/views/agenda"
)

var (
	Version    string
	CommitHash string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "task",
	Short:   "A task management CLI & TUI for your Notedown notes",
	Run:     root,
	Version: version(),
}

func root(cmd *cobra.Command, args []string) {
	cfg := loadConfig()

	// Configure logging to a file
	logFileLocation := path.Join(cfg.home, ".notedown", "logs", "task.log")
	if err := os.MkdirAll(path.Dir(logFileLocation), 0755); err != nil {
		fmt.Println("error creating log directory:", err)
		os.Exit(1)
	}
	logFile, err := os.OpenFile(logFileLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error opening log file:", err)
		os.Exit(1)
	}
	defer logFile.Close()
	slog.SetDefault(slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true})))

	client, err := notedown.NewClient(cfg.root)
	if err != nil {
		fmt.Println("error creating client:", err)
		os.Exit(1)
	}

	// Create a listener for the task client to refresh the TUI when tasks are created/updated/deleted
	// We don't need one for the daily note client because we don't need to re-render the TUI when daily notes are created/updated/deleted
	taskSub := make(chan tasks.Event)
	client.Subscribe(taskSub)
	taskListener := listeners.NewTaskListener(taskSub)

	// Create the initial model and run the program
	ctx := context.New(themes.CatpuccinMocha, context.WithListeners(taskListener))
	agenda := agenda.New(ctx, client, cfg.date)

	p := tea.NewProgram(agenda, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("error running program:", err)
		os.Exit(1)
	}
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
	home string
	root string
	date time.Time
}

func loadConfig() config {
	cfg := config{}
	cfg.root = viper.GetString("dir")
	if cfg.root == "" {
		fmt.Println("Please set NOTEDOWN_DIR environment variable to the root of your Notedown workspace")
		os.Exit(1)
	}

	// Time should always be now, but for testing purposes we allow it to be set with a hidden env var
	cfg.date = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	if t := os.Getenv("TEST_DATE"); t != "" {
		if tt, err := time.Parse("2006-01-02", t); err == nil {
			cfg.date = tt
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting user home directory:", err)
		os.Exit(1)
	}
	cfg.home = home

	return cfg
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("notedown")
	viper.BindEnv("dir")
	viper.AutomaticEnv() // read in environment variables that match
}

func version() string {
	var b strings.Builder

	if Version == "" {
		b.WriteString("dev")
	} else {
		b.WriteString(Version)
	}

	if CommitHash != "" {
		b.WriteString("-")
		b.WriteString(CommitHash)
	}

	return b.String()
}
