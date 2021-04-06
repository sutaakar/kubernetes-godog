package test

import (
	"strconv"

	"github.com/cucumber/godog"
	cloudog "github.com/sutaakar/kubernetes-godog"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	namespaceCounter := 0
	namespaceNameGenerator := func() string {
		namespaceCounter++
		return "namespace-test-" + strconv.Itoa(namespaceCounter)
	}
	cloudog.Builder().WithNamespaceNameGenerator(namespaceNameGenerator).RegisterSteps(ctx)
}
