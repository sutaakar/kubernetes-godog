package steps

import "github.com/cucumber/godog"

// KubernetesStepsBuilder builder for Kubernetes steps
type KubernetesStepsBuilder struct {
	createNamespaceListeners [](func(createdNamespace string))
	deleteNamespaceListeners [](func(deletedNamespace string))
	namespaceNameGenerator   func() string
}

// Builder returns builder for Kubernetes steps
func Builder() *KubernetesStepsBuilder {
	return &KubernetesStepsBuilder{}
}

// WithCreateNamespaceListener register listener listening for created namespace events
func (builder *KubernetesStepsBuilder) WithCreateNamespaceListener(listener func(createdNamespace string)) *KubernetesStepsBuilder {
	builder.createNamespaceListeners = append(builder.createNamespaceListeners, listener)
	return builder
}

// WithDeleteNamespaceListener register listener listening for deleted namespace events
func (builder *KubernetesStepsBuilder) WithDeleteNamespaceListener(listener func(deletedNamespace string)) *KubernetesStepsBuilder {
	builder.deleteNamespaceListeners = append(builder.deleteNamespaceListeners, listener)
	return builder
}

// WithNamespaceNameGenerator provide namespace name generator to use namespace steps with implicit names
func (builder *KubernetesStepsBuilder) WithNamespaceNameGenerator(generator func() string) *KubernetesStepsBuilder {
	builder.namespaceNameGenerator = generator
	return builder
}

// RegisterSteps register Kubernetes steps
func (builder *KubernetesStepsBuilder) RegisterSteps(ctx *godog.ScenarioContext) {
	//activeNamespace := RegisterNamespaceSteps(ctx, builder.createNamespaceListeners, builder.deleteNamespaceListeners, builder.namespaceNameGenerator)
	RegisterNamespaceSteps(ctx, builder.createNamespaceListeners, builder.deleteNamespaceListeners, builder.namespaceNameGenerator)
}
