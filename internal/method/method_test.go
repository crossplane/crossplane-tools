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

	"github.com/crossplaneio/crossplane-tools/internal/fields"
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
