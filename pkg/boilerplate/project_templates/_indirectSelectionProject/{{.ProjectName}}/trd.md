# Primary Objective

1. Create a project that will serve as an example of a gRPC client/server Application that leverages indirect selection for use in a code generation tool.  All it has to do is run, log, and provide metrics.  Don't get elaborate.  Keep it simple.

2. Use Cobra for CLI, Viper for config, and Viper automatic env for passing in environment variables.

3. Use Prometheus for metrics.  At a minimum, the example service must serve metrics for number of requests, request errors, and request durations.

4. The name of this service is "example-indirect-selection".  This name must be used as a prefix for all metrics.

5. Authentication is via https://github.com/nikogura/jwt-ssh-agent-go.

6. Application example project: ../tdoctl

7. Lint example project: ../kms

8. Note that gRPC reflection must be able to be toggled by config from the environment.

9. Follow all lint rules, including https://github.com/nikogura/namedreturns.  The golangci.yml file is syntactically correct, however the published schema does not include the custom section.  Do not alter the lint rules.  Follow them in every respect.
