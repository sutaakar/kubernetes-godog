package steps

import (
	"testing"

	"github.com/sutaakar/kubernetes-godog/pkg/core"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func Test_createNamespace(t *testing.T) {
	core.SetClient(fake.NewFakeClient())

	context := &NamespaceContext{}

	err := createNamespace(context)("test")
	if err != nil {
		t.Errorf("Error creating namespace 'test' : %v", err)
	}

	exists, err := core.IsNamespaceExists("test")
	if err != nil {
		t.Errorf("Error checking namespace 'test' : %v", err)
	}
	if !exists {
		t.Error("Namespace 'test' should exist, but it is not found")
	}

	if context.ActiveNamespace != "test" {
		t.Errorf("Active namespace contains wrong value: %s, expected to contain 'test'", context.ActiveNamespace)
	}
}

func Test_createNamespaceWithListener(t *testing.T) {
	core.SetClient(fake.NewFakeClient())

	namespaceProvidedInListener := ""
	context := &NamespaceContext{
		createNamespaceListeners: []func(createdNamespace string){func(createdNamespace string) { namespaceProvidedInListener = createdNamespace }},
	}

	err := createNamespace(context)("test")
	if err != nil {
		t.Errorf("Error creating namespace 'test' : %v", err)
	}

	exists, err := core.IsNamespaceExists("test")
	if err != nil {
		t.Errorf("Error checking namespace 'test' : %v", err)
	}
	if !exists {
		t.Error("Namespace 'test' should exist, but it is not found")
	}

	if namespaceProvidedInListener != "test" {
		t.Errorf("Created namespace listener contains wrong value: %s, expected to contain 'test'", namespaceProvidedInListener)
	}
}

func Test_createDuplicitNamespace(t *testing.T) {
	namespace := &corev1.Namespace{ObjectMeta: v1.ObjectMeta{Name: "test"}}
	core.SetClient(fake.NewFakeClient(namespace))

	context := &NamespaceContext{}

	err := createNamespace(context)("test")
	if err == nil {
		t.Error("Namespace creation should fail as there is already namespace with same name available")
	}
}

func Test_createNamespaceWithGeneratedName(t *testing.T) {
	core.SetClient(fake.NewFakeClient())

	context := &NamespaceContext{
		namespaceNameGenerator: func() string { return "generated-test" },
	}

	err := createNamespaceWithGeneratedName(context)()
	if err != nil {
		t.Errorf("Error creating namespace: %v", err)
	}

	exists, err := core.IsNamespaceExists("generated-test")
	if err != nil {
		t.Errorf("Error checking namespace 'generated-test' : %v", err)
	}
	if !exists {
		t.Error("Namespace 'generated-test' should exist, but it is not found")
	}

	if context.ActiveNamespace != "generated-test" {
		t.Errorf("Active namespace contains wrong value: %s, expected to contain 'test'", context.ActiveNamespace)
	}
}

func Test_createNamespaceWithGeneratedNameNoGenerator(t *testing.T) {
	core.SetClient(fake.NewFakeClient())

	context := &NamespaceContext{}

	err := createNamespaceWithGeneratedName(context)()
	if err == nil {
		t.Error("Should throw an error as namespace generator is not available")
	}
}

func Test_namespaceExists(t *testing.T) {
	core.SetClient(fake.NewFakeClient())

	context := &NamespaceContext{}

	err := namespaceExists("test")
	if err == nil {
		t.Error("Namespace 'test' shouldn't exist")
	}

	err = createNamespace(context)("test")
	if err != nil {
		t.Errorf("Error creating namespace 'test' : %v", err)
	}

	err = namespaceExists("test")
	if err != nil {
		t.Error("Namespace 'test' should exist")
	}
}

func Test_namespaceDoesntExist(t *testing.T) {
	core.SetClient(fake.NewFakeClient())

	context := &NamespaceContext{}

	err := namespaceDoesntExist("test")
	if err != nil {
		t.Error("Namespace 'test' shouldn't exist")
	}

	err = createNamespace(context)("test")
	if err != nil {
		t.Errorf("Error creating namespace 'test' : %v", err)
	}

	err = namespaceDoesntExist("test")
	if err == nil {
		t.Error("Namespace 'test' should exist")
	}
}

func Test_namespaceIsInState(t *testing.T) {
	namespace := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{Name: "test"},
		Status: corev1.NamespaceStatus{
			Phase: corev1.NamespaceActive,
		},
	}
	core.SetClient(fake.NewFakeClient(namespace))

	context := &NamespaceContext{
		ActiveNamespace: "test",
	}

	err := namespaceIsInState(context)(string(corev1.NamespaceActive))
	if err != nil {
		t.Errorf("Failed checking of namespace state: %v", err)
	}
}

func Test_deleteNamespace(t *testing.T) {
	namespace := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{Name: "test"},
	}
	core.SetClient(fake.NewFakeClient(namespace))

	context := &NamespaceContext{}

	err := deleteNamespace(context)("test")
	if err != nil {
		t.Errorf("Failed deleting namespace: %v", err)
	}
}

func Test_deleteNamespaceWithListener(t *testing.T) {
	namespace := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{Name: "test"},
	}
	core.SetClient(fake.NewFakeClient(namespace))

	namespaceProvidedInListener := ""
	context := &NamespaceContext{
		deleteNamespaceListeners: []func(deletedNamespace string){func(deletedNamespace string) { namespaceProvidedInListener = deletedNamespace }},
	}

	err := deleteNamespace(context)("test")
	if err != nil {
		t.Errorf("Failed deleting namespace: %v", err)
	}

	if namespaceProvidedInListener != "test" {
		t.Errorf("Deleted namespace listener contains wrong value: %s, expected to contain 'test'", namespaceProvidedInListener)
	}
}

func Test_deleteNamespaceNotExistingNamespace(t *testing.T) {
	core.SetClient(fake.NewFakeClient())

	context := &NamespaceContext{}

	err := deleteNamespace(context)("test")
	if err == nil {
		t.Error("Should throw an error as the namespace doesn't exist")
	}
}

func Test_deleteActiveNamespace(t *testing.T) {
	namespace := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{Name: "test"},
	}
	core.SetClient(fake.NewFakeClient(namespace))

	context := &NamespaceContext{
		ActiveNamespace: "test",
	}

	err := deleteActiveNamespace(context)()
	if err != nil {
		t.Errorf("Failed checking of namespace state: %v", err)
	}
}

func Test_deleteActiveNamespaceNoNamespaceActive(t *testing.T) {
	namespace := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{Name: "test"},
	}
	core.SetClient(fake.NewFakeClient(namespace))

	context := &NamespaceContext{}

	err := deleteActiveNamespace(context)()
	if err == nil {
		t.Error("Should fail as active namespace is not defined")
	}
}
