/***
Copyright 2019 Cisco Systems Inc. All rights reserved.

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

	acipolicyv1 "github.com/noironetworks/aci-containers/pkg/gbpcrd/apis/acipolicy/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeGBPSStates implements GBPSStateInterface
type FakeGBPSStates struct {
	Fake *FakeAciV1
	ns   string
}

var gbpsstatesResource = schema.GroupVersionResource{Group: "aci.aw", Version: "v1", Resource: "gbpsstates"}

var gbpsstatesKind = schema.GroupVersionKind{Group: "aci.aw", Version: "v1", Kind: "GBPSState"}

// Get takes name of the gBPSState, and returns the corresponding gBPSState object, and an error if there is any.
func (c *FakeGBPSStates) Get(ctx context.Context, name string, options v1.GetOptions) (result *acipolicyv1.GBPSState, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(gbpsstatesResource, c.ns, name), &acipolicyv1.GBPSState{})

	if obj == nil {
		return nil, err
	}
	return obj.(*acipolicyv1.GBPSState), err
}

// List takes label and field selectors, and returns the list of GBPSStates that match those selectors.
func (c *FakeGBPSStates) List(ctx context.Context, opts v1.ListOptions) (result *acipolicyv1.GBPSStateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(gbpsstatesResource, gbpsstatesKind, c.ns, opts), &acipolicyv1.GBPSStateList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &acipolicyv1.GBPSStateList{ListMeta: obj.(*acipolicyv1.GBPSStateList).ListMeta}
	for _, item := range obj.(*acipolicyv1.GBPSStateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested gBPSStates.
func (c *FakeGBPSStates) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(gbpsstatesResource, c.ns, opts))

}

// Create takes the representation of a gBPSState and creates it.  Returns the server's representation of the gBPSState, and an error, if there is any.
func (c *FakeGBPSStates) Create(ctx context.Context, gBPSState *acipolicyv1.GBPSState, opts v1.CreateOptions) (result *acipolicyv1.GBPSState, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(gbpsstatesResource, c.ns, gBPSState), &acipolicyv1.GBPSState{})

	if obj == nil {
		return nil, err
	}
	return obj.(*acipolicyv1.GBPSState), err
}

// Update takes the representation of a gBPSState and updates it. Returns the server's representation of the gBPSState, and an error, if there is any.
func (c *FakeGBPSStates) Update(ctx context.Context, gBPSState *acipolicyv1.GBPSState, opts v1.UpdateOptions) (result *acipolicyv1.GBPSState, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(gbpsstatesResource, c.ns, gBPSState), &acipolicyv1.GBPSState{})

	if obj == nil {
		return nil, err
	}
	return obj.(*acipolicyv1.GBPSState), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeGBPSStates) UpdateStatus(ctx context.Context, gBPSState *acipolicyv1.GBPSState, opts v1.UpdateOptions) (*acipolicyv1.GBPSState, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(gbpsstatesResource, "status", c.ns, gBPSState), &acipolicyv1.GBPSState{})

	if obj == nil {
		return nil, err
	}
	return obj.(*acipolicyv1.GBPSState), err
}

// Delete takes name of the gBPSState and deletes it. Returns an error if one occurs.
func (c *FakeGBPSStates) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(gbpsstatesResource, c.ns, name), &acipolicyv1.GBPSState{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeGBPSStates) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(gbpsstatesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &acipolicyv1.GBPSStateList{})
	return err
}

// Patch applies the patch and returns the patched gBPSState.
func (c *FakeGBPSStates) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *acipolicyv1.GBPSState, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(gbpsstatesResource, c.ns, name, pt, data, subresources...), &acipolicyv1.GBPSState{})

	if obj == nil {
		return nil, err
	}
	return obj.(*acipolicyv1.GBPSState), err
}
