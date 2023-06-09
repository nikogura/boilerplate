/*
	Copyright <2023> Nik Ogura <nik.ogura@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package boilerplate

import (
	"embed"
	"reflect"
	"testing"
)

func TestGetProjectFs(t *testing.T) {
	for _, tc := range []struct {
		Name    string
		Input   string
		Want    embed.FS
		WantErr bool
	}{
		{
			Name:    "Valid Cobra Project",
			Input:   "cobra",
			Want:    cobraProject,
			WantErr: false,
		},
		{
			Name:    "Invalid",
			Input:   "invalid",
			WantErr: true,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			fs, _, err := GetProjectFs(tc.Input)
			if err != nil {
				if !tc.WantErr {
					t.Errorf("unexpected err occurred: %v", err)
				}
				return
			}
			if reflect.DeepEqual(fs, embed.FS{}) {
				t.Errorf("fs nil unexpected")
			} else if !reflect.DeepEqual(fs, tc.Want) {

			}
		})
	}
}
