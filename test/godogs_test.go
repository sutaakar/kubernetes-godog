package test

import (
	"github.com/cucumber/godog"
	"github.com/sutaakar/kubernetes-godog/pkg/steps"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	steps.RegisterNamespaceSteps(ctx)
}
