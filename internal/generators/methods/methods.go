// Package methods contains methods that may be generated for a Go type.
package methods

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/negz/angryjet/internal/fields"
	"github.com/negz/angryjet/internal/generate"
)

// NewSetConditions returns a NewMethod that writes a SetConditions method for
// the supplied Object to the supplied file.
func NewSetConditions(receiver, runtime string) generate.NewMethod {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetConditions of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetConditions").Params(jen.Id("c").Op("...").Qual(runtime, "Condition")).Block(
			jen.Id(receiver).Dot(fields.NameStatus).Dot("SetConditions").Call(jen.Id("c").Op("...")),
		)
	}
}

// NewSetBindingPhase returns a NewMethod that writes a SetBindingPhase method
// for the supplied Object to the supplied file.
func NewSetBindingPhase(receiver, runtime string) generate.NewMethod {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetBindingPhase of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetBindingPhase").Params(jen.Id("p").Qual(runtime, "BindingPhase")).Block(
			jen.Id(receiver).Dot(fields.NameStatus).Dot("SetBindingPhase").Call(jen.Id("p")),
		)
	}
}

// NewGetBindingPhase returns a NewMethod that writes a GetBindingPhase method
// for the supplied Object to the supplied file.
func NewGetBindingPhase(receiver, runtime string) generate.NewMethod {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetBindingPhase of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetBindingPhase").Params().Qual(runtime, "BindingPhase").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameStatus).Dot("GetBindingPhase").Call()),
		)
	}
}

// NewSetClaimReference returns a NewMethod that writes a SetClaimReference
// method for the supplied Object to the supplied file.
func NewSetClaimReference(receiver, core string) generate.NewMethod {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetClaimReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetClaimReference").Params(jen.Id("r").Op("*").Qual(core, "ObjectReference")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("ClaimReference").Op("=").Id("r"),
		)
	}
}

// NewGetClaimReference returns a NewMethod that writes a GetClaimReference
// method for the supplied Object to the supplied file.
func NewGetClaimReference(receiver, core string) generate.NewMethod {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetClaimReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetClaimReference").Params().Op("*").Qual(core, "ObjectReference").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("ClaimReference")),
		)
	}
}

// NewSetNonPortableClassReference returns a NewMethod that writes a
// SetNonPortableClassReference method for the supplied Object to the supplied
// file.
func NewSetNonPortableClassReference(receiver, core string) generate.NewMethod {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetNonPortableClassReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetNonPortableClassReference").Params(jen.Id("r").Op("*").Qual(core, "ObjectReference")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("NonPortableClassReference").Op("=").Id("r"),
		)
	}
}

// NewGetNonPortableClassReference returns a NewMethod that writes a
// GetNonPortableClassReference method for the supplied Object to the supplied
// file.
func NewGetNonPortableClassReference(receiver, core string) generate.NewMethod {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetNonPortableClassReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetNonPortableClassReference").Params().Op("*").Qual(core, "ObjectReference").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("NonPortableClassReference")),
		)
	}
}

// NewSetWriteConnectionSecretToReference returns a NewMethod that writes a
// SetWriteConnectionSecretToReference method for the supplied Object to the
// supplied file.
func NewSetWriteConnectionSecretToReference(receiver, core string) generate.NewMethod {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetWriteConnectionSecretToReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetWriteConnectionSecretToReference").Params(jen.Id("r").Qual(core, "LocalObjectReference")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("WriteConnectionSecretToReference").Op("=").Id("r"),
		)
	}
}

// NewGetWriteConnectionSecretToReference returns a NewMethod that writes a
// GetWriteConnectionSecretToReference method for the supplied Object to the
// supplied file.
func NewGetWriteConnectionSecretToReference(receiver, core string) generate.NewMethod {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetWriteConnectionSecretToReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetWriteConnectionSecretToReference").Params().Qual(core, "LocalObjectReference").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("WriteConnectionSecretToReference")),
		)
	}
}

// NewSetReclaimPolicy returns a NewMethod that writes a SetReclaimPolicy method
// for the supplied Object to the supplied file.
func NewSetReclaimPolicy(receiver, core string) generate.NewMethod {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetReclaimPolicy of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetReclaimPolicy").Params(jen.Id("r").Qual(core, "ReclaimPolicy")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("ReclaimPolicy").Op("=").Id("r"),
		)
	}
}

// NewGetReclaimPolicy returns a NewMethod that writes a GetReclaimPolicy method
// for the supplied Object to the supplied file.
func NewGetReclaimPolicy(receiver, runtime string) generate.NewMethod {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetReclaimPolicy of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetReclaimPolicy").Params().Qual(runtime, "ReclaimPolicy").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("ReclaimPolicy")),
		)
	}
}
