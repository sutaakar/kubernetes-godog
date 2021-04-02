package steps

import (
	"errors"
	"fmt"

	"github.com/cucumber/godog"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}

		if err := create(ns); err != nil {
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
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}

		if err := create(ns); err != nil {
			return err
		}

		for _, listener := range context.createNamespaceListeners {
			listener(namespaceName)
		}

		return nil
	}
}

func namespaceExists(namespaceName string) error {
	if _, err := getNamespace(namespaceName); err != nil {
		return err
	}
	return nil
}

func namespaceDoesntExist(namespaceName string) error {
	if _, err := getNamespace(namespaceName); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	return fmt.Errorf("Namespace %s found", namespaceName)
}

func namespaceIsInState(context *NamespaceContext) func(string) error {
	return func(namespacePhase string) error {
		if ns, err := getNamespace(context.ActiveNamespace); err != nil {
			return err
		} else if ns.Status.Phase != corev1.NamespacePhase(namespacePhase) {
			return fmt.Errorf("Expected namespace phase %s, but got %s", namespacePhase, ns.Status.Phase)
		}
		return nil
	}
}

func deleteNamespace(context *NamespaceContext) func(string) error {
	return func(namespaceName string) error {
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}

		if err := delete(ns); err != nil {
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
		return deleteNamespace(context)(context.ActiveNamespace)
	}
}

// ### Utility methods

func getNamespace(namespaceName string) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{}
	err := get(types.NamespacedName{Name: namespaceName}, ns)

	return ns, err
}
