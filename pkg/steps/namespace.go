package steps

import (
	"fmt"

	"github.com/cucumber/godog"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// RegisterNamespaceSteps registers all steps related to namespace operations
func RegisterNamespaceSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^create namespace ([a-z0-9-]+)$`, createNamespace)
	ctx.Step(`^namespace ([a-z0-9-]+) exists$`, namespaceExists)
	ctx.Step(`^namespace ([a-z0-9-]+) doesn't exist$`, namespaceDoesntExist)
	ctx.Step(`^namespace is in state ([a-zA-Z]+)$`, namespaceIsInState)
	ctx.Step(`^delete namespace ([a-z0-9-]+)$`, deleteNamespace)
	ctx.Step(`^delete namespace$`, deleteActiveNamespace)
}

func createNamespace(namespaceName string) error {
	Context.ActiveNamespace = namespaceName
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}
	return create(ns)
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

func namespaceIsInState(namespacePhase string) error {
	if ns, err := getNamespace(Context.ActiveNamespace); err != nil {
		return err
	} else if ns.Status.Phase != corev1.NamespacePhase(namespacePhase) {
		return fmt.Errorf("Expected namespace phase %s, but got %s", namespacePhase, ns.Status.Phase)
	}
	return nil
}

func deleteNamespace(namespaceName string) error {
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}
	return delete(ns)
}

func deleteActiveNamespace() error {
	return deleteNamespace(Context.ActiveNamespace)
}

// ### Utility methods

func getNamespace(namespaceName string) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{}
	err := get(types.NamespacedName{Name: namespaceName}, ns)

	return ns, err
}
