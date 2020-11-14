# kubernetes-godog

This project contains definition and implementation of Godog steps related to various Kubernetes objects.
These steps can be easily added to existing test suite by passing ScenarioContext to one of available functions for step registration.
Step implementations use standard kubeconfig approach to connect to Kubernetes cluster (using env variable to point to config file or config file in expected location)

Currently there are specified just namespace steps, more steps will come.
