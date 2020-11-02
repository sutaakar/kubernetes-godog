package test

import (
	"strconv"

	"github.com/cucumber/godog"
	"github.com/sutaakar/kubernetes-godog/pkg/steps"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	namespaceCounter := 0
	namespaceNameGenerator := func() string {
		namespaceCounter++
		return "namespace-test-" + strconv.Itoa(namespaceCounter)
	}
	steps.RegisterGeneratedNamespaceSteps(ctx, namespaceNameGenerator)
}
