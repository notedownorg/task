// Copyright 2024 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

// sandboxCmd represents the sandbox command
var sandboxCmd = &cobra.Command{
	Use:    "sandbox",
	Hidden: true,
	Short:  "Run using sandbox workspace",
	Run: func(cmd *cobra.Command, args []string) {

		// Assumes being ran from git root
		rootDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		// TODO: We should completely randomly generate the workspace contents

		sandboxPath := path.Join(rootDir, "sandbox")
		os.Setenv("NOTEDOWN_DIR", sandboxPath)
		root(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(sandboxCmd)

}
