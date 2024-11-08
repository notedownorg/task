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

package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"sync"

	cp "github.com/otiai10/copy"
	"github.com/spf13/pflag"
)

var (
	generateGifs = pflag.Bool("generate-gifs", false, "generate gifs for each tape")
)

func main() {
	pflag.Parse()

	dirs, err := dirs("features")
	handleErr(err)

	var wg sync.WaitGroup
	errors := make([]error, 0)
	for _, dir := range dirs {
		wg.Add(1)
		go func(dir string) {
			defer wg.Done()
			start := fmt.Sprintf("./%s/start", dir)
			final := fmt.Sprintf("./%s/final", dir)
			tape := fmt.Sprintf("./%s/scenario.tape", dir)
			gif := fmt.Sprintf("./%s/scenario.gif", dir)

			// Setup the workspace by copying the start directory to a fresh final directory
			slog.Info("setting up workspace", "start", start, "final", final)
			if err := os.RemoveAll(final); err != nil {
				errors = append(errors, err)
				return
			}
			if err := cp.Copy(start, final); err != nil {
				errors = append(errors, err)
				return
			}

			// Run the tape
			slog.Info("running tape", "tape", tape)
			cmd := exec.Command("vhs", tape)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				errors = append(errors, err)
				return
			}

			// If generateGifs is false, remove the gifs
			if !*generateGifs {
				if err := os.Remove(gif); err != nil {
					errors = append(errors, err)
					return
				}
			}
		}(dir)
	}

	wg.Wait()
	if len(errors) > 0 {
		for _, err := range errors {
			slog.Error(err.Error())
		}
		os.Exit(1)
	}
}

func dirs(root string) ([]string, error) {
	dirs, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0)
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		out = append(out, fmt.Sprintf("%s/%s", root, dir.Name()))
	}
	return out, nil
}

func handleErr(err error) {
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
