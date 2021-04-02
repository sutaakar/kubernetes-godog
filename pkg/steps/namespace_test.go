package steps

import (
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func Test_CreateNamespace(t *testing.T) {
	k8sClient = fake.NewFakeClient()

	err := CreateNamespace("test")
	if err != nil {
		t.Errorf("Error creating namespace 'test' : %v", err)
	}

	exists, err := IsNamespaceExists("test")
	if err != nil {
		t.Errorf("Error checking namespace 'test' : %v", err)
	}
	if !exists {
		t.Error("Namespace 'test' should exist, but it is not found")
	}
}

func Test_IsNamespaceExists(t *testing.T) {
	k8sClient = fake.NewFakeClient()

	exists, err := IsNamespaceExists("test")
	if err != nil {
		t.Errorf("Error checking namespace 'test' : %v", err)
	}
	if exists {
		t.Error("Namespace 'test' shouldn't exist, but it is found")
	}

	err = CreateNamespace("test")
	if err != nil {
		t.Errorf("Error creating namespace 'test' : %v", err)
	}

	exists, err = IsNamespaceExists("test")
	if err != nil {
		t.Errorf("Error checking namespace 'test' : %v", err)
	}
	if !exists {
		t.Error("Namespace 'test' should exist, but it is not found")
	}
}

func Test_DeleteNamespace(t *testing.T) {
	k8sClient = fake.NewFakeClient()

	err := CreateNamespace("test")
	if err != nil {
		t.Errorf("Error creating namespace 'test' : %v", err)
	}

	exists, err := IsNamespaceExists("test")
	if err != nil {
		t.Errorf("Error checking namespace 'test' : %v", err)
	}
	if !exists {
		t.Error("Namespace 'test' should exist, but it is not found")
	}

	err = DeleteNamespace("test")
	if err != nil {
		t.Errorf("Error deleting namespace 'test' : %v", err)
	}

	exists, err = IsNamespaceExists("test")
	if err != nil {
		t.Errorf("Error checking namespace 'test' : %v", err)
	}
	if exists {
		t.Error("Namespace 'test' shouldn't exist, but it is found")
	}
}

func Test_DeleteNotExistingNamespace(t *testing.T) {
	k8sClient = fake.NewFakeClient()

	err := DeleteNamespace("test")
	if err == nil {
		t.Errorf("Expected error when deleting not existing namespace")
	}
}
