/*
	Copyright <2023> Nik Ogura <nik.ogura@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package boilerplate

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/spf13/afero"
	"io"
	"path/filepath"
	"strings"
	"text/template"
)

type FilePath struct {
	Path      string
	Name      string
	TemplPath string
	TemplName string
	IsDir     bool
}

type TmplWriter struct {
	OutFs     afero.Fs
	TemplFs   embed.FS
	FilePaths []FilePath
	ProjDir   string
	TmplVals  map[string]any
}

func NewTmplWriter(outFs afero.Fs, projType string, vals map[string]any) (TmplWriter, error) {
	fs, dirName, err := GetProjectFs(projType)
	if err != nil {
		return TmplWriter{}, fmt.Errorf("fs error: %w", err)
	}

	w := TmplWriter{
		OutFs:    outFs,
		TemplFs:  fs,
		ProjDir:  dirName,
		TmplVals: vals}

	w.FilePaths, err = w.GetFilePaths(dirName)
	if err != nil {
		return w, fmt.Errorf("failed to walk filepath from root(%s): %w", ".", err)
	}
	return w, nil
}

func (w TmplWriter) BuildProject(destDir string) error {
	//fp := w.FilePaths
	err := w.ResolveAllPathTemplates()
	if err != nil {
		return err
	}

	w.fixGoModTemplPaths()

	err = w.CreateAllFilePathsAtRoot(destDir)
	if err != nil {
		return err
	}

	err = w.WriteAllDestFileTemplateData(destDir)
	if err != nil {
		return err
	}

	return nil
}

func (w TmplWriter) ResolveAllPathTemplates() error {
	for i := range w.FilePaths {
		fp := w.FilePaths[i]
		buf, err := w.ResolveTemplateVars(fp.Path)
		if err != nil {
			return fmt.Errorf("path resolution failure: path=%s, err=%w", fp.Path, err)
		} else {
			path := strings.Replace(buf.String(), w.ProjDir, "", 1)
			if path[0] == '/' {
				path = path[1:]
			}
			w.FilePaths[i].TemplPath = path
		}

		buf, err = w.ResolveTemplateVars(fp.Name)
		if err != nil {
			return fmt.Errorf("name resolution failure: path=%s, err=%w", fp.Name, err)
		} else {
			name := strings.Replace(buf.String(), w.ProjDir, "", 1)
			if name[0] == '/' {
				name = name[1:]
			}
			w.FilePaths[i].TemplName = name
		}
	}

	return nil
}

func (w TmplWriter) ResolveTemplateVars(str string) (*bytes.Buffer, error) {
	tmpl, err := template.New("tmplWriter").Parse(str)
	if err != nil {
		return nil, fmt.Errorf("path parsing error: %w", err)
	}

	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, w.TmplVals)
	if err != nil {
		return nil, fmt.Errorf("failed to exec template from string(%s): %w", str, err)
	}

	return buf, nil
}

func (w TmplWriter) CreateAllFilePathsAtRoot(root string) error {
	for _, fp := range w.FilePaths {
		err := w.CreatePath(root, fp.TemplPath, fp.IsDir)
		if err != nil {
			return fmt.Errorf("failed to create file(%s): err(%w)", fp.TemplPath, err)
		}
	}
	return nil
}

func (w TmplWriter) CreatePath(root, file string, isDir bool) error {
	path := fmt.Sprintf("%s/%s", root, file)
	if isDir {
		// Create the directory itself
		err := w.OutFs.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("directory creation failed: %w", err)
		}
	} else {
		// Create parent directories for files
		err := w.OutFs.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return fmt.Errorf("directory creation failed: %w", err)
		}
	}
	return nil
}

func (w TmplWriter) WriteAllDestFileTemplateData(destDir string) error {
	for _, fp := range w.FilePaths {
		if fp.IsDir {
			continue
		}

		err := w.WriteFileTemplateData(fp, destDir)
		if err != nil {
			return fmt.Errorf("writing file(%s) failed: %w", fp.TemplName, err)
		}
	}

	return nil
}

func (w TmplWriter) WriteFileTemplateData(fp FilePath, destDir string) error {
	path := fmt.Sprintf("%s/%s", destDir, fp.TemplPath)
	file, err := w.OutFs.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file(%s): %w", path, err)
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil { //nolint:staticcheck // close error handling not needed
			// Log or handle close error if needed
		}
	}()

	buf, err := w.ResolveFileTemplateData(fp)
	if err != nil {
		return fmt.Errorf("failed to exec template: %w", err)
	}

	buf, err = w.removeBuildExclusions(buf)
	if err != nil {
		return fmt.Errorf("failed to remove build exclusions from file(%s): %w", fp.TemplName, err)
	}

	n, err := file.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("cannot write file(%s) bytes: %w", path, err)
	} else if n != buf.Len() {
		return fmt.Errorf("wrong number of bytes written: exp(%d) act(%d)", buf.Len(), n)
	}

	return nil
}

func (w TmplWriter) ResolveFileTemplateData(fp FilePath) (*bytes.Buffer, error) {
	// Read the original file data not the parsed template path
	file, err := w.TemplFs.Open(fp.Path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file(%s): err(%w)", fp.Path, err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("cannot read file data(%s): err(%w)", fp.Path, err)
	}

	buf, err := w.ResolveTemplateVars(string(data))
	if err != nil {
		return nil, fmt.Errorf("cannot execute template with file(%s): %w", fp.Path, err)
	}

	return buf, nil
}

func (w TmplWriter) GetFilePaths(root string) ([]FilePath, error) {
	var fp []FilePath

	entries, err := w.TemplFs.ReadDir(root)
	if err != nil {
		return fp, fmt.Errorf("failed to read embedded files at dir(%s): %w", root, err)
	}

	for _, e := range entries {
		cpath := fmt.Sprintf("%s/%s", root, e.Name())
		if e.IsDir() {
			// Add directory entry
			fp = append(fp, FilePath{
				Path:  cpath,
				Name:  e.Name(),
				IsDir: true,
			})

			children, childErr := w.GetFilePaths(cpath)
			if childErr != nil {
				return nil, fmt.Errorf("failed to collected children at root path(%s): %w", cpath, childErr)
			}
			fp = append(fp, children...)
		} else {
			fp = append(fp, FilePath{
				Path:  cpath,
				Name:  e.Name(),
				IsDir: false,
			})
		}
	}

	return fp, nil
}

func (w TmplWriter) fixGoModTemplPaths() {
	for i := range w.FilePaths {
		fp := w.FilePaths[i]
		if fp.TemplName == "go.mod_" || fp.TemplName == "go.sum_" {
			fp.TemplPath = fp.TemplPath[:len(fp.TemplPath)-1]
			fp.TemplName = fp.TemplName[:len(fp.TemplName)-1]
		}
		w.FilePaths[i] = fp
	}
}

func (w TmplWriter) removeBuildExclusions(buf *bytes.Buffer) (*bytes.Buffer, error) {
	replBuf, err := w.ResolveTemplateVars("// +build exclude {{.ProjectName}}\n")
	if err != nil {
		return nil, fmt.Errorf("removing build exclusions cannot execute template: %w", err)
	}

	return bytes.NewBufferString(strings.ReplaceAll(buf.String(), replBuf.String(), "")), nil
}
