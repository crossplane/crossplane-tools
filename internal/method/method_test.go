/*
Copyright 2019 The Crossplane Authors.

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

package method

import (
	"fmt"
	"go/types"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/google/go-cmp/cmp"
)

type MockObject struct {
	types.Object

	Named string
}

func (o MockObject) Name() string {
	return o.Named
}

func TestNewSetConditions(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// SetConditions of this Type.
func (t *Type) SetConditions(c ...runtime.Condition) {
	t.Status.SetConditions(c...)
}
`
	f := jen.NewFile("pkg")
	NewSetConditions("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetConditions(): -want, +got\n%s", diff)
	}
}

func TestNewGetCondition(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// GetCondition of this Type.
func (t *Type) GetCondition(ct runtime.ConditionType) runtime.Condition {
	return t.Status.GetCondition(ct)
}
`
	f := jen.NewFile("pkg")
	NewGetCondition("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetCondition(): -want, +got\n%s", diff)
	}
}

func TestNewSetResourceReference(t *testing.T) {
	want := `package pkg

import core "example.org/core"

// SetResourceReference of this Type.
func (t *Type) SetResourceReference(r *core.ObjectReference) {
	t.Spec.ResourceReference = r
}
`
	f := jen.NewFile("pkg")
	NewSetResourceReference("t", "example.org/core")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetResourceReference(): -want, +got\n%s", diff)
	}
}

func TestNewGetResourceReference(t *testing.T) {
	want := `package pkg

import core "example.org/core"

// GetResourceReference of this Type.
func (t *Type) GetResourceReference() *core.ObjectReference {
	return t.Spec.ResourceReference
}
`
	f := jen.NewFile("pkg")
	NewGetResourceReference("t", "example.org/core")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetResourceReference(): -want, +got\n%s", diff)
	}
}

func TestNewSetProviderConfigReference(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// SetProviderConfigReference of this Type.
func (t *Type) SetProviderConfigReference(r *runtime.Reference) {
	t.Spec.ProviderConfigReference = r
}
`
	f := jen.NewFile("pkg")
	NewSetProviderConfigReference("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetProviderConfigReference(): -want, +got\n%s", diff)
	}
}

func TestNewGetProviderConfigReference(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// GetProviderConfigReference of this Type.
func (t *Type) GetProviderConfigReference() *runtime.Reference {
	return t.Spec.ProviderConfigReference
}
`
	f := jen.NewFile("pkg")
	NewGetProviderConfigReference("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetProviderReference(): -want, +got\n%s", diff)
	}
}

func TestNewSetWriteConnectionSecretToReference(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// SetWriteConnectionSecretToReference of this Type.
func (t *Type) SetWriteConnectionSecretToReference(r *runtime.SecretReference) {
	t.Spec.WriteConnectionSecretToReference = r
}
`
	f := jen.NewFile("pkg")
	NewSetWriteConnectionSecretToReference("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetWriteConnectionSecretToReference(): -want, +got\n%s", diff)
	}
}

func TestNewGetWriteConnectionSecretToReference(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// GetWriteConnectionSecretToReference of this Type.
func (t *Type) GetWriteConnectionSecretToReference() *runtime.SecretReference {
	return t.Spec.WriteConnectionSecretToReference
}
`
	f := jen.NewFile("pkg")
	NewGetWriteConnectionSecretToReference("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetWriteConnectionSecretToLocalReference(): -want, +got\n%s", diff)
	}
}

func TestNewLocalSetWriteConnectionSecretToReference(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// SetWriteConnectionSecretToReference of this Type.
func (t *Type) SetWriteConnectionSecretToReference(r *runtime.LocalSecretReference) {
	t.Spec.WriteConnectionSecretToReference = r
}
`
	f := jen.NewFile("pkg")
	NewLocalSetWriteConnectionSecretToReference("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetWriteConnectionSecretToLocalReference(): -want, +got\n%s", diff)
	}
}

func TestNewLocalGetWriteConnectionSecretToReference(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// GetWriteConnectionSecretToReference of this Type.
func (t *Type) GetWriteConnectionSecretToReference() *runtime.LocalSecretReference {
	return t.Spec.WriteConnectionSecretToReference
}
`
	f := jen.NewFile("pkg")
	NewLocalGetWriteConnectionSecretToReference("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetWriteConnectionSecretToReference(): -want, +got\n%s", diff)
	}
}

func TestNewSetDeletionPolicy(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// SetDeletionPolicy of this Type.
func (t *Type) SetDeletionPolicy(r runtime.DeletionPolicy) {
	t.Spec.DeletionPolicy = r
}
`
	f := jen.NewFile("pkg")
	NewSetDeletionPolicy("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetDeletionPolicy(): -want, +got\n%s", diff)
	}
}

func TestNewGetDeletionPolicy(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// GetDeletionPolicy of this Type.
func (t *Type) GetDeletionPolicy() runtime.DeletionPolicy {
	return t.Spec.DeletionPolicy
}
`
	f := jen.NewFile("pkg")
	NewGetDeletionPolicy("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetDeletionPolicy(): -want, +got\n%s", diff)
	}
}

func TestNewSetUsers(t *testing.T) {
	want := `package pkg

// SetUsers of this Type.
func (t *Type) SetUsers(i int64) {
	t.Status.Users = i
}
`
	f := jen.NewFile("pkg")
	NewSetUsers("t")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetUsers(): -want, +got\n%s", diff)
	}
}

func TestNewGetUsers(t *testing.T) {
	want := `package pkg

// GetUsers of this Type.
func (t *Type) GetUsers() int64 {
	return t.Status.Users
}
`
	f := jen.NewFile("pkg")
	NewGetUsers("t")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetUsers(): -want, +got\n%s", diff)
	}
}

func TestNewManagedGetItems(t *testing.T) {
	want := `package pkg

import resource "example.org/resource"

// GetItems of this Type.
func (t *Type) GetItems() []resource.Managed {
	items := make([]resource.Managed, len(t.Items))
	for i := range t.Items {
		items[i] = &t.Items[i]
	}
	return items
}
`
	f := jen.NewFile("pkg")
	NewManagedGetItems("t", "example.org/resource")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewManagedGetItems(): -want, +got\n%s", diff)
	}
}

func TestNewSetRootProviderConfigReference(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// SetProviderConfigReference of this Type.
func (t *Type) SetProviderConfigReference(r runtime.Reference) {
	t.ProviderConfigReference = r
}
`
	f := jen.NewFile("pkg")
	NewSetRootProviderConfigReference("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetRootProviderConfigReference(): -want, +got\n%s", diff)
	}
}

func TestNewGetRootProviderConfigReference(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// GetProviderConfigReference of this Type.
func (t *Type) GetProviderConfigReference() runtime.Reference {
	return t.ProviderConfigReference
}
`
	f := jen.NewFile("pkg")
	NewGetRootProviderConfigReference("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetRootProviderReference(): -want, +got\n%s", diff)
	}
}

func TestNewSetRootResourceReference(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// SetResourceReference of this Type.
func (t *Type) SetResourceReference(r runtime.TypedReference) {
	t.ResourceReference = r
}
`
	f := jen.NewFile("pkg")
	NewSetRootResourceReference("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetRootResourceReference(): -want, +got\n%s", diff)
	}
}

func TestNewGetRootResourceReference(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// GetResourceReference of this Type.
func (t *Type) GetResourceReference() runtime.TypedReference {
	return t.ResourceReference
}
`
	f := jen.NewFile("pkg")
	NewGetRootResourceReference("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetRootProviderReference(): -want, +got\n%s", diff)
	}
}

func TestNewProviderConfigUsageGetItems(t *testing.T) {
	want := `package pkg

import resource "example.org/resource"

// GetItems of this Type.
func (t *Type) GetItems() []resource.ProviderConfigUsage {
	items := make([]resource.ProviderConfigUsage, len(t.Items))
	for i := range t.Items {
		items[i] = &t.Items[i]
	}
	return items
}
`
	f := jen.NewFile("pkg")
	NewProviderConfigUsageGetItems("t", "example.org/resource")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewProviderConfigUsageGetItems(): -want, +got\n%s", diff)
	}
}
