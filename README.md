# kubernetes-godog

This project contains definition and implementation of Godog steps related to various Kubernetes objects.
These steps can be easily added to existing test suite by passing ScenarioContext to the builder.
Step implementations use standard kubeconfig approach to connect to Kubernetes cluster (using env variable to point to config file or config file in expected location)

Currently there are specified just namespace steps, more steps will come.

## Usage

Register steps in InitializeScenario:
```go
func InitializeScenario(ctx *godog.ScenarioContext) {
	steps.Builder().RegisterSteps(ctx)
}
```

Use registered steps in your feature files:
```gherkin
  Scenario: My complex scenario
    When create namespace dedicated-namespace

    Then namespace dedicated-namespace exists
```
