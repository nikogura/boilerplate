# Primary Objective

1. Create a project that will serve as an example of a Single-Page Application runnnig out of an embedded filesystem for use in a code generation tool.  All it has to do is run, log, and provide metrics.  Don't get elaborate.  Keep it simple.

2. Use Cobra for CLI, Viper for config, and Viper automatic env for passing in environment variables.

3. Use Prometheus for metrics.  At a minimum, the example service must serve metrics for number of requests, request errors, and request durations.

4. The name of this service is "{{.ProjectName}}".  This name must be used as a prefix for all metrics.

5. Authentication is via OIDC

6. Application example: ../depctl

7. Template example: ../boilerplate  Take care when writing the templates, because the github action templates need to be double-escaped.  Also the go.mod and go.sum files need to be sourced verbatim.  The goal is for a user, once they generate code via the boilerplate tool will only need to run 'go build && ./<project name>'.
