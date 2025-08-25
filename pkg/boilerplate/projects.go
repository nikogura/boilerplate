/*
	Copyright <2023> Nik Ogura <nik.ogura@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package boilerplate

import (
	"embed"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"os"
)

const (
	VERSION               = "3.6.0"
	CobraProjectType      = "cobra"
	HeadlessServiceType   = "headless-service"
	SPAProjectType        = "spa"
	IndirectSelectionType = "indirect-selection"
)

//go:embed project_templates/_cobraProject/*
var cobraProject embed.FS

//go:embed project_templates/_headlessServiceProject/*
var headlessServiceProject embed.FS

//go:embed project_templates/_spaProject/*
var spaProject embed.FS

//go:embed all:project_templates/_indirectSelectionProject
var indirectSelectionProject embed.FS

// GetProjectFs  Gets the embedded file system for the project of this type.
func GetProjectFs(projType string) (embed.FS, string, error) {
	switch projType {
	case CobraProjectType:
		return cobraProject, "project_templates/_cobraProject", nil
	case HeadlessServiceType:
		return headlessServiceProject, "project_templates/_headlessServiceProject", nil
	case SPAProjectType:
		return spaProject, "project_templates/_spaProject", nil
	case IndirectSelectionType:
		return indirectSelectionProject, "project_templates/_indirectSelectionProject", nil
	}

	return embed.FS{}, "", fmt.Errorf("failed to detect embedded package: %s", projType)
}

// ValidProjectTypes  Lists the valid project types.
func ValidProjectTypes() []string {
	return []string{
		CobraProjectType,
		HeadlessServiceType,
		SPAProjectType,
		IndirectSelectionType,
	}
}

// IsValidProjectType  Returns true or false depending on whether the project is a supported type.
func IsValidProjectType(v string) bool {
	switch v {
	case CobraProjectType:
		return true
	case HeadlessServiceType:
		return true
	case SPAProjectType:
		return true
	case IndirectSelectionType:
		return true
	}
	return false
}

// promptForParamsWithRetry handles the retry loop for parameter collection.
func promptForParamsWithRetry[T PromptValues](data T, promptFunc func(T, io.Reader) error) T {
	for {
		err := promptFunc(data, os.Stdin)
		if err != nil {
			fmt.Print(color.RedString("%s\n", err))
		} else {
			return data
		}
	}
}

func PromptsForProject(proj string) (data PromptValues, err error) {
	switch proj {
	case CobraProjectType:
		return promptForParamsWithRetry(&CobraCliToolParams{}, CobraCliToolParamsFromPrompts), nil

	case HeadlessServiceType:
		return promptForParamsWithRetry(&HeadlessServiceParams{}, HeadlessServiceParamsFromPrompts), nil

	case SPAProjectType:
		return promptForParamsWithRetry(NewSPAParams(), SPAParamsFromPrompts), nil

	case IndirectSelectionType:
		return promptForParamsWithRetry(&IndirectSelectionParams{}, IndirectSelectionParamsFromPrompts), nil

	default:
		log.Fatalf("unknown or unhandled project type. options are %s", ValidProjectTypes())
	}

	return data, err
}
