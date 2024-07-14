/*
Copyright 2024 The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "github.com/openshift/node-feature-discovery/api/nfd/v1alpha1"
)

// FakeNodeFeatureGroups implements NodeFeatureGroupInterface
type FakeNodeFeatureGroups struct {
	Fake *FakeNfdV1alpha1
	ns   string
}

var nodefeaturegroupsResource = v1alpha1.SchemeGroupVersion.WithResource("nodefeaturegroups")

var nodefeaturegroupsKind = v1alpha1.SchemeGroupVersion.WithKind("NodeFeatureGroup")

// Get takes name of the nodeFeatureGroup, and returns the corresponding nodeFeatureGroup object, and an error if there is any.
func (c *FakeNodeFeatureGroups) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.NodeFeatureGroup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(nodefeaturegroupsResource, c.ns, name), &v1alpha1.NodeFeatureGroup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NodeFeatureGroup), err
}

// List takes label and field selectors, and returns the list of NodeFeatureGroups that match those selectors.
func (c *FakeNodeFeatureGroups) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.NodeFeatureGroupList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(nodefeaturegroupsResource, nodefeaturegroupsKind, c.ns, opts), &v1alpha1.NodeFeatureGroupList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.NodeFeatureGroupList{ListMeta: obj.(*v1alpha1.NodeFeatureGroupList).ListMeta}
	for _, item := range obj.(*v1alpha1.NodeFeatureGroupList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested nodeFeatureGroups.
func (c *FakeNodeFeatureGroups) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(nodefeaturegroupsResource, c.ns, opts))

}

// Create takes the representation of a nodeFeatureGroup and creates it.  Returns the server's representation of the nodeFeatureGroup, and an error, if there is any.
func (c *FakeNodeFeatureGroups) Create(ctx context.Context, nodeFeatureGroup *v1alpha1.NodeFeatureGroup, opts v1.CreateOptions) (result *v1alpha1.NodeFeatureGroup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(nodefeaturegroupsResource, c.ns, nodeFeatureGroup), &v1alpha1.NodeFeatureGroup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NodeFeatureGroup), err
}

// Update takes the representation of a nodeFeatureGroup and updates it. Returns the server's representation of the nodeFeatureGroup, and an error, if there is any.
func (c *FakeNodeFeatureGroups) Update(ctx context.Context, nodeFeatureGroup *v1alpha1.NodeFeatureGroup, opts v1.UpdateOptions) (result *v1alpha1.NodeFeatureGroup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(nodefeaturegroupsResource, c.ns, nodeFeatureGroup), &v1alpha1.NodeFeatureGroup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NodeFeatureGroup), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeNodeFeatureGroups) UpdateStatus(ctx context.Context, nodeFeatureGroup *v1alpha1.NodeFeatureGroup, opts v1.UpdateOptions) (*v1alpha1.NodeFeatureGroup, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(nodefeaturegroupsResource, "status", c.ns, nodeFeatureGroup), &v1alpha1.NodeFeatureGroup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NodeFeatureGroup), err
}

// Delete takes name of the nodeFeatureGroup and deletes it. Returns an error if one occurs.
func (c *FakeNodeFeatureGroups) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(nodefeaturegroupsResource, c.ns, name, opts), &v1alpha1.NodeFeatureGroup{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeNodeFeatureGroups) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(nodefeaturegroupsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.NodeFeatureGroupList{})
	return err
}

// Patch applies the patch and returns the patched nodeFeatureGroup.
func (c *FakeNodeFeatureGroups) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.NodeFeatureGroup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(nodefeaturegroupsResource, c.ns, name, pt, data, subresources...), &v1alpha1.NodeFeatureGroup{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NodeFeatureGroup), err
}
