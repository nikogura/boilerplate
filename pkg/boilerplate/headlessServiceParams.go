/*
	Copyright <2022> Nik Ogura <nik.ogura@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package boilerplate

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"strings"
)

type HeadlessServiceParams struct {
	ProjectName       string `json:"ProjectName"`
	ProjectPackage    string `json:"ProjectPackage"`
	EnvPrefix         string `json:"EnvPrefix"`
	ProjectShortDesc  string `json:"ProjectShortDesc"`
	ProjectLongDesc   string `json:"ProjectLongDesc"`
	MaintainerName    string `json:"MaintainerName"`
	MaintainerEmail   string `json:"MaintainerEmail"`
	GolangVersion     string `json:"GolangVersion"`
	DbtRepo           string `json:"DbtRepo"`
	ProjectVersion    string `json:"ProjectVersion"`
	DefaultServerPort string `json:"DefaultServerPort"`
	ServerShortDesc   string `json:"ServerShortDesc"`
	ServerLongDesc    string `json:"ServerLongDesc"`
	OwnerName         string `json:"OwnerName"`
	OwnerEmail        string `json:"OwnerEmail"`
}

func (hsp *HeadlessServiceParams) Values() map[ParamPrompt]*string {
	return map[ParamPrompt]*string{
		GoVersion:           &hsp.GolangVersion,
		DockerRegistry:      nil,
		DockerProject:       nil,
		ProjName:            &hsp.ProjectName,
		ProjPkgName:         &hsp.ProjectPackage,
		ProjEnvPrefix:       &hsp.EnvPrefix,
		ProjShortDesc:       &hsp.ProjectShortDesc,
		ProjLongDesc:        &hsp.ProjectLongDesc,
		ProjMaintainerName:  &hsp.MaintainerName,
		ProjMaintainerEmail: &hsp.MaintainerEmail,
		DbtRepo:             &hsp.DbtRepo,
		ProjectVersion:      &hsp.ProjectVersion,
		ServerDefPort:       &hsp.DefaultServerPort,
		ServerShortDesc:     &hsp.ServerShortDesc,
		ServerLongDesc:      &hsp.ServerLongDesc,
		OwnerName:           &hsp.OwnerName,
		OwnerEmail:          &hsp.OwnerEmail,
	}
}

func (hsp *HeadlessServiceParams) AsMap() (output map[string]any, err error) {
	data, err := json.Marshal(&hsp)
	if err != nil {
		err = errors.Wrapf(err, "failed to marshal params object")
		return output, err
	}

	output = make(map[string]any)
	err = json.Unmarshal(data, &output)
	if err != nil {
		err = errors.Wrapf(err, "failed to unmarshal data just marshalled")
		return output, err
	}

	// Add a Go package-safe version of ProjectName
	output["ProjectPackageName"] = strings.ReplaceAll(hsp.ProjectName, "-", "")

	return output, err
}

func GetHeadlessServiceParamsPromptMessaging() map[ParamPrompt]Prompt {
	prompts := commonPromptMessaging()

	// Add headless service specific prompts
	prompts[ProjEnvPrefix] = Prompt{
		PromptMsg:    "Enter environment variable prefix for your service.",
		InputFailMsg: "failed to read environment prefix",
		Validations:  envPrefix,
		DefaultValue: "SERVICE",
	}

	prompts[ServerDefPort] = Prompt{
		PromptMsg:    "Enter default metrics port.",
		InputFailMsg: "failed to read default metrics port",
		Validations:  portValidation,
		DefaultValue: "8080",
	}

	prompts[OwnerName] = Prompt{
		PromptMsg:    "Enter the owner/organization name.",
		InputFailMsg: "failed to read owner name",
		DefaultValue: "you@example.com",
	}

	prompts[OwnerEmail] = Prompt{
		PromptMsg:    "Enter the owner/organization email address.",
		InputFailMsg: "failed to read owner email address",
		Validations:  emailValidation,
		DefaultValue: "code@example.com",
	}

	return prompts
}

func HeadlessServiceParamsFromPrompts(params *HeadlessServiceParams, r io.Reader) (err error) {
	prompts := GetHeadlessServiceParamsPromptMessaging()
	err = paramsFromPrompts(r, prompts, params)
	if err != nil {
		return err
	}

	// Auto-populate server descriptions from project descriptions
	if params.ServerShortDesc == "" {
		params.ServerShortDesc = params.ProjectShortDesc
	}
	if params.ServerLongDesc == "" {
		params.ServerLongDesc = params.ProjectLongDesc
	}

	return err
}
