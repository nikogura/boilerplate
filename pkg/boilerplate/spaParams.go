/*
	Copyright <2022> Nik Ogura <nik.ogura@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package boilerplate

import (
	"io"
	"strings"
)

// SPAParams holds all the parameters for SPA project templates.
type SPAParams struct {
	ProjectName            *string
	ProjectPackage         *string
	ProjectPackageName     *string
	ProjectShortDesc       *string
	ProjectLongDesc        *string
	ProjectVersion         *string
	GolangVersion          *string
	ProjectMaintainerName  *string
	ProjectMaintainerEmail *string
	DbtRepo                *string
}

// NewSPAParams creates a new SPAParams with initialized string pointers.
func NewSPAParams() *SPAParams {
	return &SPAParams{
		ProjectName:            new(string),
		ProjectPackage:         new(string),
		ProjectPackageName:     new(string),
		ProjectShortDesc:       new(string),
		ProjectLongDesc:        new(string),
		ProjectVersion:         new(string),
		GolangVersion:          new(string),
		ProjectMaintainerName:  new(string),
		ProjectMaintainerEmail: new(string),
		DbtRepo:                new(string),
	}
}

// Values returns a map of all parameter values.
func (p *SPAParams) Values() map[ParamPrompt]*string {
	return map[ParamPrompt]*string{
		GoVersion:           p.GolangVersion,
		DockerRegistry:      nil,
		DockerProject:       nil,
		ProjName:            p.ProjectName,
		ProjPkgName:         p.ProjectPackage,
		ProjEnvPrefix:       nil,
		ProjShortDesc:       p.ProjectShortDesc,
		ProjLongDesc:        p.ProjectLongDesc,
		ProjMaintainerName:  p.ProjectMaintainerName,
		ProjMaintainerEmail: p.ProjectMaintainerEmail,
		ServerDefPort:       nil,
		ServerShortDesc:     nil,
		ServerLongDesc:      nil,
		OwnerName:           nil,
		OwnerEmail:          nil,
		DbtRepo:             p.DbtRepo,
		ProjectVersion:      p.ProjectVersion,
	}
}

// AsMap converts the parameters to a map for template processing.
func (p *SPAParams) AsMap() (data map[string]any, err error) {
	data = map[string]any{
		"ProjectName":            *p.ProjectName,
		"ProjectPackage":         *p.ProjectPackage,
		"ProjectShortDesc":       *p.ProjectShortDesc,
		"ProjectLongDesc":        *p.ProjectLongDesc,
		"ProjectVersion":         *p.ProjectVersion,
		"GolangVersion":          *p.GolangVersion,
		"ProjectMaintainerName":  *p.ProjectMaintainerName,
		"ProjectMaintainerEmail": *p.ProjectMaintainerEmail,
		"DbtRepo":                *p.DbtRepo,
	}

	// Add a Go package-safe version of ProjectName
	data["ProjectPackageName"] = strings.ReplaceAll(*p.ProjectName, "-", "")

	// Add uppercase version for environment variables
	envPrefix := strings.ToUpper(strings.ReplaceAll(*p.ProjectName, "-", "_"))
	data["ProjectEnvPrefix"] = envPrefix

	return data, err
}

// SPAParamsFromPrompts populates SPA parameters from user prompts.
func SPAParamsFromPrompts(p *SPAParams, r io.Reader) (err error) {
	prompts := commonPromptMessaging()

	return paramsFromPrompts(r, prompts, p)
}
