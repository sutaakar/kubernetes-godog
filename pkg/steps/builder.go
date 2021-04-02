package steps

import "github.com/cucumber/godog"

// KubernetesStepsBuilder builder for Kubernetes steps
type KubernetesStepsBuilder struct {
	namespaceContext NamespaceContext
}

// Builder returns builder for Kubernetes steps
func Builder() *KubernetesStepsBuilder {
	return &KubernetesStepsBuilder{}
}

// WithCreateNamespaceListener register listener listening for created namespace events
func (builder *KubernetesStepsBuilder) WithCreateNamespaceListener(listener func(createdNamespace string)) *KubernetesStepsBuilder {
	builder.namespaceContext.createNamespaceListeners = append(builder.namespaceContext.createNamespaceListeners, listener)
	return builder
}

// WithDeleteNamespaceListener register listener listening for deleted namespace events
func (builder *KubernetesStepsBuilder) WithDeleteNamespaceListener(listener func(deletedNamespace string)) *KubernetesStepsBuilder {
	builder.namespaceContext.deleteNamespaceListeners = append(builder.namespaceContext.deleteNamespaceListeners, listener)
	return builder
}

// WithNamespaceNameGenerator provide namespace name generator to use namespace steps with implicit names
func (builder *KubernetesStepsBuilder) WithNamespaceNameGenerator(generator func() string) *KubernetesStepsBuilder {
	builder.namespaceContext.namespaceNameGenerator = generator
	return builder
}

// RegisterSteps register Kubernetes steps
func (builder *KubernetesStepsBuilder) RegisterSteps(ctx *godog.ScenarioContext) {
	RegisterNamespaceSteps(ctx, &builder.namespaceContext)
}
