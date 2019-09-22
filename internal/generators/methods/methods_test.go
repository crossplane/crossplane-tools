package methods

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

func TestNewSetNonPortableClassReference(t *testing.T) {
	want := `package pkg

import core "example.org/core"

// SetNonPortableClassReference of this Type.
func (t *Type) SetNonPortableClassReference(r *core.ObjectReference) {
	t.Spec.NonPortableClassReference = r
}
`
	f := jen.NewFile("pkg")
	NewSetNonPortableClassReference("t", "example.org/core")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetNonPortableClassReference(): -want, +got\n%s", diff)
	}
}

func TestNewGetNonPortableClassReference(t *testing.T) {
	want := `package pkg

import core "example.org/core"

// GetNonPortableClassReference of this Type.
func (t *Type) GetNonPortableClassReference() *core.ObjectReference {
	return t.Spec.NonPortableClassReference
}
`
	f := jen.NewFile("pkg")
	NewGetNonPortableClassReference("t", "example.org/core")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetNonPortableClassReference(): -want, +got\n%s", diff)
	}
}

func TestNewSetWriteConnectionSecretToReference(t *testing.T) {
	want := `package pkg

import core "example.org/core"

// SetWriteConnectionSecretToReference of this Type.
func (t *Type) SetWriteConnectionSecretToReference(r core.LocalObjectReference) {
	t.Spec.WriteConnectionSecretToReference = r
}
`
	f := jen.NewFile("pkg")
	NewSetWriteConnectionSecretToReference("t", "example.org/core")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewSetWriteConnectionSecretToReference(): -want, +got\n%s", diff)
	}
}

func TestNewGetWriteConnectionSecretToReference(t *testing.T) {
	want := `package pkg

import core "example.org/core"

// GetWriteConnectionSecretToReference of this Type.
func (t *Type) GetWriteConnectionSecretToReference() core.LocalObjectReference {
	return t.Spec.WriteConnectionSecretToReference
}
`
	f := jen.NewFile("pkg")
	NewGetWriteConnectionSecretToReference("t", "example.org/core")(f, MockObject{Named: "Type"})
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
	NewSetReclaimPolicy("t", "example.org/runtime")(f, MockObject{Named: "Type"})
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
	NewGetReclaimPolicy("t", "example.org/runtime")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewGetReclaimPolicy(): -want, +got\n%s", diff)
	}
}
