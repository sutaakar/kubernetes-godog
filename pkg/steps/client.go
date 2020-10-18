package steps

import (
	"fmt"
	"os"

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
			fmt.Printf("failed to create client %v", err)
			os.Exit(1)
		}
	}
	return k8sClient
}
