/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apihelper

import (
	api "k8s.io/api/core/v1"
	k8sclient "k8s.io/client-go/kubernetes"
)

// APIHelpers represents a set of API helpers for Kubernetes
type APIHelpers interface {
	// GetClient returns a client
	GetClient() (*k8sclient.Clientset, error)

	// GetNode returns the Kubernetes node on which this container is running.
	GetNode(*k8sclient.Clientset, string) (*api.Node, error)

	// UpdateNode updates the node via the API server using a client.
	UpdateNode(*k8sclient.Clientset, *api.Node) error
}
