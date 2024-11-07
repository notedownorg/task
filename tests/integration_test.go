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

package tests

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	cp "github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "addtask",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Copy workspace into a temp directory at a deterministic location
			// This must be deterministic as it has to be hardcoded in the tape
			tmp := fmt.Sprintf("testdata/tmp/%s", tt.name)
			golden := fmt.Sprintf("testdata/golden/%s", tt.name)

			// Ensure the dir is clean and copy the starting workspace
			assert.NoError(t, os.RemoveAll(tmp))
			assert.NoError(t, cp.Copy("testdata/workspace", tmp))

			// Run the tape
			cmd := exec.Command("vhs", fmt.Sprintf("testdata/%s.tape", tt.name), "-o", "testdata/tmp/addtask.gif")
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Errorf("Output: %s, Error: %s", output, err)
			}

			// Walk testdata/tmp and check that the testdata/golden equivalent matches
			assert.NoError(t, filepath.Walk(tmp, func(path string, info fs.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				tmpContents, err := os.ReadFile(path)
				assert.NoError(t, err)
				goldenContents, err := os.ReadFile(filepath.Join(golden, path[len(tmp):]))
				assert.NoError(t, err)
				assert.Equal(t, string(goldenContents), string(tmpContents))
				return nil // we check errors as we go
			}))

			// Check that the golden directory has no extra files
			assert.NoError(t, filepath.Walk(golden, func(path string, info fs.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				_, err = os.Stat(filepath.Join(tmp, path[len(golden):]))
				assert.NoError(t, err)
				return nil // we check errors as we go
			}))
		})
	}

}
