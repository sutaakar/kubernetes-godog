package steps

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// CreateNamespace creates namespace
func CreateNamespace(namespaceName string) error {
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}
	return create(ns)
}

// IsNamespaceExists returns true is namespace exists
func IsNamespaceExists(namespaceName string) (bool, error) {
	if _, err := GetNamespace(namespaceName); err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetNamespace returns namespace
func GetNamespace(namespaceName string) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{}
	err := get(types.NamespacedName{Name: namespaceName}, ns)

	return ns, err
}

// DeleteNamespace deletes namespace
func DeleteNamespace(namespaceName string) error {
	ns, err := GetNamespace(namespaceName)
	if err != nil {
		return err
	}
	return delete(ns)
}
