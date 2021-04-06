package steps

import (
	"errors"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/sutaakar/kubernetes-godog/pkg/core"
	corev1 "k8s.io/api/core/v1"
)

// RegisterNamespaceSteps registers all steps related to namespace operations
func RegisterNamespaceSteps(ctx *godog.ScenarioContext, createNamespaceListeners [](func(createdNamespace string)), deleteNamespaceListeners [](func(deletedNamespace string)), namespaceNameGenerator func() string) *string {
	activeNamespace := ""
	ctx.Step(`^create namespace ([a-z0-9-]+)$`, createNamespace(&activeNamespace, createNamespaceListeners))
	ctx.Step(`^create namespace$`, createNamespaceWithGeneratedName(&activeNamespace, createNamespaceListeners, namespaceNameGenerator))
	ctx.Step(`^namespace ([a-z0-9-]+) exists$`, namespaceExists)
	ctx.Step(`^namespace ([a-z0-9-]+) doesn't exist$`, namespaceDoesntExist)
	ctx.Step(`^namespace is in ([a-zA-Z]+) state$`, namespaceIsInState(&activeNamespace))
	ctx.Step(`^delete namespace ([a-z0-9-]+)$`, deleteNamespace(deleteNamespaceListeners))
	ctx.Step(`^delete namespace$`, deleteActiveNamespace(&activeNamespace, deleteNamespaceListeners))
	return &activeNamespace
}

func createNamespace(activeNamespace *string, createNamespaceListeners [](func(createdNamespace string))) func(string) error {
	return func(namespaceName string) error {
		*activeNamespace = namespaceName
		if err := core.CreateNamespace(namespaceName); err != nil {
			return err
		}

		for _, listener := range createNamespaceListeners {
			listener(namespaceName)
		}

		return nil
	}
}

func createNamespaceWithGeneratedName(activeNamespace *string, createNamespaceListeners [](func(createdNamespace string)), namespaceNameGenerator func() string) func() error {
	return func() error {
		if namespaceNameGenerator == nil {
			return errors.New("Namespace name generator not specified. Please provide a namespace name generator")
		}

		namespaceName := namespaceNameGenerator()
		*activeNamespace = namespaceName
		if err := core.CreateNamespace(namespaceName); err != nil {
			return err
		}

		for _, listener := range createNamespaceListeners {
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

func namespaceIsInState(activeNamespace *string) func(string) error {
	return func(namespacePhase string) error {
		if ns, err := core.GetNamespace(*activeNamespace); err != nil {
			return err
		} else if ns.Status.Phase != corev1.NamespacePhase(namespacePhase) {
			return fmt.Errorf("Expected namespace phase %s, but got %s", namespacePhase, ns.Status.Phase)
		}
		return nil
	}
}

func deleteNamespace(deleteNamespaceListeners [](func(deletedNamespace string))) func(string) error {
	return func(namespaceName string) error {
		if err := core.DeleteNamespace(namespaceName); err != nil {
			return err
		}

		for _, listener := range deleteNamespaceListeners {
			listener(namespaceName)
		}

		return nil
	}
}

func deleteActiveNamespace(activeNamespace *string, deleteNamespaceListeners [](func(deletedNamespace string))) func() error {
	return func() error {
		if len(*activeNamespace) == 0 {
			return fmt.Errorf("Active namespace not defined")
		}
		return deleteNamespace(deleteNamespaceListeners)(*activeNamespace)
	}
}
