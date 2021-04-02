package steps

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var k8sClient client.Client

func getClient() client.Client {
	if k8sClient == nil {
		config := config.GetConfigOrDie()
		// Adjust config values to reduce throttling
		config.QPS = 100.0
		config.Burst = 150.0

		var err error
		k8sClient, err = client.New(config, client.Options{})
		if err != nil {
			panic(fmt.Sprintf("failed to create client %v", err))
		}
	}
	return k8sClient
}

// Retrieves object based on a given key identifier
func get(key client.ObjectKey, object runtime.Object) error {
	return getClient().Get(context.TODO(), key, object)
}

// Creates object using options if provided
func create(object runtime.Object, opts ...client.CreateOption) error {
	return getClient().Create(context.TODO(), object, opts...)
}

// Deletes object using options if provided
func delete(object runtime.Object, opts ...client.DeleteOption) error {
	return getClient().Delete(context.TODO(), object, opts...)
}
