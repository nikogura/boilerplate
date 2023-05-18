# Boilerplate

[![Current Release](https://img.shields.io/github/release/nikogura/boilerplate.svg)](https://img.shields.io/github/release/nikogura/boilerplate.svg)

[![Circle CI](https://circleci.com/gh/nikogura/boilerplate.svg?style=shield)](https://circleci.com/gh/nikogura/boilerplate)

[![Go Report Card](https://goreportcard.com/badge/github.com/nikogura/boilerplate)](https://goreportcard.com/report/github.com/nikogura/boilerplate)

[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/nikogura/boilerplate/pkg/boilerplate)

[![Coverage Status](https://codecov.io/gh/nikogura/boilerplate/branch/master/graph/badge.svg)](https://codecov.io/gh/nikogura/boilerplate)

This tool uses a templated file system to help generate templated projects.
Each of the folders in this directory contain a layout for a specific type of
project to generate.  Each folder name starts with an underscore (_) to prevent go tools from treating them as actual golang project files.

Within a given project, all items are templatized; folders and files.

## Project Types
### [Cobra](pkg/boilerplate/project_templates/_cobraProject)
This project is used to generate tools using the [cobra](https://github.com/spf13/cobra) command line framework.

## Adding a new Project
### Make a project folder
First step is to creat a new "projects" folder in the [project_templates](pkg/boilerplate/project_templates) directory. Under this
created directory you can create any number of templated file structures that will become the basic of your
new destination project.

For example, under a project you might create both a templated service and web GUI application which require
separate templating schemes.


NB: Your directory name needs to start with an underscore ("_").  This will ensure the golang tools ignore it.  If you don't follow this rule, things like `go mod` will throw errors on the template syntax.

NB: Only single folder projects have been attempted at the time of this writing.

### Add project to [projects.go](pkg/boilerplate/projects.go)
Create a go:embed FS to hold your project structure
```shell script
go:embed project_templates/_cobraProject/*
var myNewProject embed.FS
```

Add the project to each function in this file.

### Add new prompt types
If adding new template variables, they should be added to the [prompt.go](../prompt.go) file. This
includes the prompt questions as well as any validations to perform on a given answer.

### Make params structure
Create a struct that holds all of the variables your application requires to run

```
type DockerParams struct {
   	DockerRegistry    string
   	DockerProject     string
   	ProjectName       string
    ...
```

### Done
After this build the binary via `go build`, and your new project will be available for generation at the top level of the application

## Use with Gomason

The Boilerplate tool was designed for use with [https://github.com/nikogura/gomason](gomason) tool.  You certainly don't have to use `gomason`, but the example cobra project included here was intended to be used that way.

You can build the project simply by running `go build` or `go install`, but if you're going to build, test, sign, and upload, you might want to run `gomason publish -vsl` instead.  (publish verbose, skip tests, operate off of what's here on the local disk in this directory).  

Gomason's other options and features can be investigated at [https://github.com/nikogura/gomason](the gomason project page).