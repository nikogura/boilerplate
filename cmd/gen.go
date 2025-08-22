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
	"bufio"
	"fmt"
	"github.com/nikogura/boilerplate/pkg/boilerplate"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
)

var projectType string //nolint:gochecknoglobals // cobra command flag
var destDir string     //nolint:gochecknoglobals // cobra command flag

// promptForProjectType prompts the user to select a project type from available options.
func promptForProjectType() string {
	validTypes := boilerplate.ValidProjectTypes()

	fmt.Println("Available project types:")
	for i, pType := range validTypes {
		fmt.Printf("  %d. %s\n", i+1, pType)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Please select a project type (number): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Please enter a valid number.")
			continue
		}

		if choice < 1 || choice > len(validTypes) {
			fmt.Printf("Please enter a number between 1 and %d.\n", len(validTypes))
			continue
		}

		return validTypes[choice-1]
	}
}

// genCmd represents the create command.
var genCmd = &cobra.Command{ //nolint:gochecknoglobals // cobra command definition
	Use:   "gen",
	Short: "Creates a new code project based on an embedded template.",
	Long: `
Creates a new code project based on an embedded template.

When run, it will ask you for a project name, and description.  It will also ask you for the name and email address of the author.  This information is used to generate all the boilerplate files and such.

Creates a skeleton project based on one of several project types.  Current supported types are:

cobra   -   A project based on the excellent Cobra CLI framework.
headless-service    -   A project implementing a standalone headless service useful for implementing APIs and the like.
spa     -   A project based on React, designed to be built as a self-contained single page application.

Each project is set up so it can be built, and provides CI workflows for both DBT tools as well as Github actions.

You can specify the project type on the command line, or be prompted.  If you omit the destination directory, it defaults to the CWD.

	`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Determine project type
		if projectType == "" {
			if len(args) > 0 {
				projectType = args[0]
			} else {
				// Prompt user for project type selection
				projectType = promptForProjectType()
			}
		}

		if destDir == "" {
			destDir, err = os.Getwd()
			if err != nil {
				log.Fatalf("failed to determine CWD: %v", err)
			}
		}

		if !boilerplate.IsValidProjectType(projectType) {
			log.Fatalf("invalid project type: %q. Valid project types are: %s", projectType, boilerplate.ValidProjectTypes())
		}

		fmt.Printf("Creating new project of type %q\n", projectType)

		prompts, err := boilerplate.PromptsForProject(projectType)
		if err != nil {
			log.Fatalf("failed to get prompts for project type %s: %v", projectType, err)
		}

		datamap, err := prompts.AsMap()
		if err != nil {
			log.Fatalf("failed to export params as map: %v", err)
		}

		wr, err := boilerplate.NewTmplWriter(afero.NewOsFs(), projectType, datamap)
		if err != nil {
			log.Fatalf("failed to create template writer: %v", err)
		}

		err = wr.BuildProject(destDir)
		if err != nil {
			log.Fatalf("failed to create templated project: %v", err)
		}

		fmt.Printf("New project created in ./%s\n", datamap["ProjectName"])
	},
}

func init() { //nolint:gochecknoinits // cobra command registration
	RootCmd.AddCommand(genCmd)
	genCmd.Flags().StringVarP(&projectType, "type", "t", "", "Project Type (if not specified, you'll be prompted to select)")
	genCmd.Flags().StringVarP(&destDir, "dest-dir", "d", "", "Destination Directory (Defaults to CWD)")
}
