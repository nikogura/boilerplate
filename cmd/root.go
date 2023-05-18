// Copyright © 2023 Nik Ogura <nik.ogura@gmail.com>
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

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "boilerplate",
	Short: "Code Generation Tool",
	Long: `
Code Generation Tool.

You can create a project any old way, including writing every file from scratch. Sometimes, however, you have better things to do with your time, and just want to get something working.  Other times you might have less sophisticated users that need to be able to get *something* up and running in a standard way without having to know all the ins and outs of your build system.

The "boilerplate" tool is intended to create a minimally usable codebase so you can get on to the real task of solving problems.
`,
}

// Execute - execute the command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
