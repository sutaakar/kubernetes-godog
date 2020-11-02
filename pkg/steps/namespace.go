package steps

import (
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
	ActiveNamespace string
}

// RegisterNamespaceSteps registers all steps related to namespace operations
func RegisterNamespaceSteps(ctx *godog.ScenarioContext) *NamespaceContext {
	context := &NamespaceContext{}
	ctx.Step(`^create namespace ([a-z0-9-]+)$`, createNamespace(context))
	ctx.Step(`^namespace ([a-z0-9-]+) exists$`, namespaceExists)
	ctx.Step(`^namespace ([a-z0-9-]+) doesn't exist$`, namespaceDoesntExist)
	ctx.Step(`^namespace is in state ([a-zA-Z]+)$`, namespaceIsInState(context))
	ctx.Step(`^delete namespace ([a-z0-9-]+)$`, deleteNamespace)
	ctx.Step(`^delete namespace$`, deleteActiveNamespace(context))
	return context
}

// RegisterGeneratedNamespaceSteps registers all steps related to generated namespace operations
func RegisterGeneratedNamespaceSteps(ctx *godog.ScenarioContext, generateNamespaceName func() string) *NamespaceContext {
	// Contains all usual namespace steps
	context := RegisterNamespaceSteps(ctx)
	ctx.Step(`^create namespace$`, createNamespaceWithGeneratedName(context, generateNamespaceName))
	return context
}

func createNamespace(context *NamespaceContext) func(string) error {
	return func(namespaceName string) error {
		context.ActiveNamespace = namespaceName
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}
		return create(ns)
	}
}

func createNamespaceWithGeneratedName(context *NamespaceContext, generateNamespaceName func() string) func() error {
	return func() error {
		namespaceName := generateNamespaceName()
		context.ActiveNamespace = namespaceName
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}
		return create(ns)
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

func deleteNamespace(namespaceName string) error {
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}
	return delete(ns)
}

func deleteActiveNamespace(context *NamespaceContext) func() error {
	return func() error {
		return deleteNamespace(context.ActiveNamespace)
	}
}

// ### Utility methods

func getNamespace(namespaceName string) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{}
	err := get(types.NamespacedName{Name: namespaceName}, ns)

	return ns, err
}
