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

	"github.com/crossplane/crossplane-tools/internal/fields"
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

func TestNewSetBindingPhase(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// SetBindingPhase of this Type.
func (t *Type) SetBindingPhase(p runtime.BindingPhase) {
	t.Status.SetBindingPhase(p)
}
`
	f := jen.NewFile("pkg")
	NewSetBindingPhase("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetBindingPhase(): -want, +got\n%s", diff)
	}
}

func TestNewGetBindingPhase(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// GetBindingPhase of this Type.
func (t *Type) GetBindingPhase() runtime.BindingPhase {
	return t.Status.GetBindingPhase()
}
`
	f := jen.NewFile("pkg")
	NewGetBindingPhase("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetBindingPhase(): -want, +got\n%s", diff)
	}
}

func TestNewSetClaimReference(t *testing.T) {
	want := `package pkg

import core "example.org/core"

// SetClaimReference of this Type.
func (t *Type) SetClaimReference(r *core.ObjectReference) {
	t.Spec.ClaimReference = r
}
`
	f := jen.NewFile("pkg")
	NewSetClaimReference("t", "example.org/core")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetClaimReference(): -want, +got\n%s", diff)
	}
}

func TestNewGetClaimReference(t *testing.T) {
	want := `package pkg

import core "example.org/core"

// GetClaimReference of this Type.
func (t *Type) GetClaimReference() *core.ObjectReference {
	return t.Spec.ClaimReference
}
`
	f := jen.NewFile("pkg")
	NewGetClaimReference("t", "example.org/core")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetClaimReference(): -want, +got\n%s", diff)
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

func TestNewSetClassSelector(t *testing.T) {
	want := `package pkg

import meta "example.org/meta"

// SetClassSelector of this Type.
func (t *Type) SetClassSelector(s *meta.LabelSelector) {
	t.Spec.ClassSelector = s
}
`
	f := jen.NewFile("pkg")
	NewSetClassSelector("t", "example.org/meta")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetClassSelector(): -want, +got\n%s", diff)
	}
}

func TestNewGetClassSelector(t *testing.T) {
	want := `package pkg

import meta "example.org/meta"

// GetClassSelector of this Type.
func (t *Type) GetClassSelector() *meta.LabelSelector {
	return t.Spec.ClassSelector
}
`
	f := jen.NewFile("pkg")
	NewGetClassSelector("t", "example.org/meta")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetClassSelector(): -want, +got\n%s", diff)
	}
}

func TestNewSetClassReference(t *testing.T) {
	want := `package pkg

import core "example.org/core"

// SetClassReference of this Type.
func (t *Type) SetClassReference(r *core.ObjectReference) {
	t.Spec.ClassReference = r
}
`
	f := jen.NewFile("pkg")
	NewSetClassReference("t", "example.org/core")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetClassReference(): -want, +got\n%s", diff)
	}
}

func TestNewGetClassReference(t *testing.T) {
	want := `package pkg

import core "example.org/core"

// GetClassReference of this Type.
func (t *Type) GetClassReference() *core.ObjectReference {
	return t.Spec.ClassReference
}
`
	f := jen.NewFile("pkg")
	NewGetClassReference("t", "example.org/core")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetClassReference(): -want, +got\n%s", diff)
	}
}

func TestNewSetProviderReference(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// SetProviderReference of this Type.
func (t *Type) SetProviderReference(r runtime.Reference) {
	t.Spec.ProviderReference = r
}
`
	f := jen.NewFile("pkg")
	NewSetProviderReference("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetProviderReference(): -want, +got\n%s", diff)
	}
}

func TestNewGetProviderReference(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// GetProviderReference of this Type.
func (t *Type) GetProviderReference() runtime.Reference {
	return t.Spec.ProviderReference
}
`
	f := jen.NewFile("pkg")
	NewGetProviderReference("t", "example.org/runtime")(f, MockObject{Named: "Type"})
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

func TestNewSetReclaimPolicy(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// SetReclaimPolicy of this Type.
func (t *Type) SetReclaimPolicy(r runtime.ReclaimPolicy) {
	t.Spec.ReclaimPolicy = r
}
`
	f := jen.NewFile("pkg")
	NewSetReclaimPolicy("t", "example.org/runtime", fields.NameSpec)(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetReclaimPolicy(): -want, +got\n%s", diff)
	}
}

func TestNewGetReclaimPolicy(t *testing.T) {
	want := `package pkg

import runtime "example.org/runtime"

// GetReclaimPolicy of this Type.
func (t *Type) GetReclaimPolicy() runtime.ReclaimPolicy {
	return t.Spec.ReclaimPolicy
}
`
	f := jen.NewFile("pkg")
	NewGetReclaimPolicy("t", "example.org/runtime", fields.NameSpec)(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetReclaimPolicy(): -want, +got\n%s", diff)
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

func TestNewClaimGetItems(t *testing.T) {
	want := `package pkg

import resource "example.org/resource"

// GetItems of this Type.
func (t *Type) GetItems() []resource.Claim {
	items := make([]resource.Claim, len(t.Items))
	for i := range t.Items {
		items[i] = &t.Items[i]
	}
	return items
}
`
	f := jen.NewFile("pkg")
	NewClaimGetItems("t", "example.org/resource")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewClaimGetItems(): -want, +got\n%s", diff)
	}
}

func TestNewClassGetItems(t *testing.T) {
	want := `package pkg

import resource "example.org/resource"

// GetItems of this Type.
func (t *Type) GetItems() []resource.Class {
	items := make([]resource.Class, len(t.Items))
	for i := range t.Items {
		items[i] = &t.Items[i]
	}
	return items
}
`
	f := jen.NewFile("pkg")
	NewClassGetItems("t", "example.org/resource")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewClassGetItems(): -want, +got\n%s", diff)
	}
}
