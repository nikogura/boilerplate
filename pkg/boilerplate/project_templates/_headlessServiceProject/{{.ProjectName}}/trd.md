# Primary Objective

1. Provide an example of a headless service (one that does not serve any external clients beyond prometheus) for use in a code generation tool.  All it has to do is run, log, and provide metrics.  Don't get elaborate.  Keep it simple.

2. Use Cobra for CLI, Viper for config, and Viper automatic env for passing in environment variables.

3. Use Prometheus for metrics.  At a minimum, the example service must serve metrics for number of requests, request errors, and request durations.

4. The name of this service is "{{.ProjectName}}".  This name must be used as a prefix for all metrics.