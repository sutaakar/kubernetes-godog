package steps

import (
	"errors"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/sutaakar/kubernetes-godog/pkg/core"
	corev1 "k8s.io/api/core/v1"
)

// NamespaceContext carries contextual information related to namespace step execution and results
type NamespaceContext struct {
	// ActiveNamespace stores latest created or used namespace
	ActiveNamespace          string
	createNamespaceListeners [](func(createdNamespace string))
	deleteNamespaceListeners [](func(deletedNamespace string))
	namespaceNameGenerator   func() string
}

// RegisterNamespaceSteps registers all steps related to namespace operations
func RegisterNamespaceSteps(ctx *godog.ScenarioContext, context *NamespaceContext) {
	ctx.Step(`^create namespace ([a-z0-9-]+)$`, createNamespace(context))
	ctx.Step(`^create namespace$`, createNamespaceWithGeneratedName(context))
	ctx.Step(`^namespace ([a-z0-9-]+) exists$`, namespaceExists)
	ctx.Step(`^namespace ([a-z0-9-]+) doesn't exist$`, namespaceDoesntExist)
	ctx.Step(`^namespace is in ([a-zA-Z]+) state$`, namespaceIsInState(context))
	ctx.Step(`^delete namespace ([a-z0-9-]+)$`, deleteNamespace(context))
	ctx.Step(`^delete namespace$`, deleteActiveNamespace(context))
}

func createNamespace(context *NamespaceContext) func(string) error {
	return func(namespaceName string) error {
		context.ActiveNamespace = namespaceName
		if err := core.CreateNamespace(namespaceName); err != nil {
			return err
		}

		for _, listener := range context.createNamespaceListeners {
			listener(namespaceName)
		}

		return nil
	}
}

func createNamespaceWithGeneratedName(context *NamespaceContext) func() error {
	return func() error {
		if context.namespaceNameGenerator == nil {
			return errors.New("Namespace name generator not specified. Please provide a namespace name generator")
		}

		namespaceName := context.namespaceNameGenerator()
		context.ActiveNamespace = namespaceName
		if err := core.CreateNamespace(namespaceName); err != nil {
			return err
		}

		for _, listener := range context.createNamespaceListeners {
			listener(namespaceName)
		}

		return nil
	}
}

func namespaceExists(namespaceName string) error {
	exists, err := core.IsNamespaceExists(namespaceName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Namespace %s doesn't exist", namespaceName)
	}
	return nil
}

func namespaceDoesntExist(namespaceName string) error {
	exists, err := core.IsNamespaceExists(namespaceName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Namespace %s exists", namespaceName)
	}
	return nil
}

func namespaceIsInState(context *NamespaceContext) func(string) error {
	return func(namespacePhase string) error {
		if ns, err := core.GetNamespace(context.ActiveNamespace); err != nil {
			return err
		} else if ns.Status.Phase != corev1.NamespacePhase(namespacePhase) {
			return fmt.Errorf("Expected namespace phase %s, but got %s", namespacePhase, ns.Status.Phase)
		}
		return nil
	}
}

func deleteNamespace(context *NamespaceContext) func(string) error {
	return func(namespaceName string) error {
		if err := core.DeleteNamespace(namespaceName); err != nil {
			return err
		}

		for _, listener := range context.deleteNamespaceListeners {
			listener(namespaceName)
		}

		return nil
	}
}

func deleteActiveNamespace(context *NamespaceContext) func() error {
	return func() error {
		if len(context.ActiveNamespace) == 0 {
			return fmt.Errorf("Active namespace not defined")
		}
		return deleteNamespace(context)(context.ActiveNamespace)
	}
}
