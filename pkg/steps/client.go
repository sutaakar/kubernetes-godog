package steps

import (
	"context"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var k8sClient client.Client

func init() {
	if k8sClient == nil {
		config := config.GetConfigOrDie()
		// Adjust config values to reduce throttling
		config.QPS = 100.0
		config.Burst = 150.0

		var err error
		k8sClient, err = client.New(config, client.Options{})
		if err != nil {
			fmt.Printf("failed to create client %v", err)
			os.Exit(1)
		}
	}
}

// Retrieves object based on a given key identifier
func get(key client.ObjectKey, object runtime.Object) error {
	return k8sClient.Get(context.TODO(), key, object)
}

// Creates object using options if provided
func create(object runtime.Object, opts ...client.CreateOption) error {
	return k8sClient.Create(context.TODO(), object, opts...)
}

// Deletes object using options if provided
func delete(object runtime.Object, opts ...client.DeleteOption) error {
	return k8sClient.Delete(context.TODO(), object, opts...)
}
