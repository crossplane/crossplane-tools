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

// Package method contains methods that may be generated for a Go type.
package method

import (
	"go/token"
	"go/types"
	"sort"

	"github.com/dave/jennifer/jen"

	"github.com/crossplane/crossplane-tools/internal/fields"
)

// New is a function that adds a method on the supplied object in the
// supplied file.
type New func(f *jen.File, o types.Object)

// A Set is a map of method names to the New functions that produce
// them.
type Set map[string]New

// Write the method Set for the supplied Object to the supplied file. Methods
// are filtered by the supplied Filter.
func (s Set) Write(f *jen.File, o types.Object, mf Filter) {
	names := make([]string, 0, len(s))
	for name := range s {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if mf(o, name) {
			continue
		}
		s[name](f, o)
	}
}

// A Filter is a function that determines whether a method should be written for
// the supplied object. It returns true if the method should be filtered.
type Filter func(o types.Object, methodName string) bool

// DefinedOutside returns a MethodFilter that returns true if the supplied
// object has a method with the supplied name that is not defined in the
// supplied filename. The object's filename is determined using the supplied
// FileSet.
func DefinedOutside(fs *token.FileSet, filename string) Filter {
	return func(o types.Object, name string) bool {
		s := types.NewMethodSet(types.NewPointer(o.Type()))
		for i := 0; i < s.Len(); i++ {
			mo := s.At(i).Obj()
			if mo.Name() != name {
				continue
			}
			if fs.Position(mo.Pos()).Filename != filename {
				return true
			}
		}
		return false
	}
}

// NewSetConditions returns a NewMethod that writes a SetConditions method for
// the supplied Object to the supplied file.
func NewSetConditions(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetConditions of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetConditions").Params(jen.Id("c").Op("...").Qual(runtime, "Condition")).Block(
			jen.Id(receiver).Dot(fields.NameStatus).Dot("SetConditions").Call(jen.Id("c").Op("...")),
		)
	}
}

// NewGetCondition returns a NewMethod that writes a GetCondition method for
// the supplied Object to the supplied file.
func NewGetCondition(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetCondition of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetCondition").Params(jen.Id("ct").Qual(runtime, "ConditionType")).Qual(runtime, "Condition").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameStatus).Dot("GetCondition").Call(jen.Id("ct"))),
		)
	}
}

// NewSetBindingPhase returns a NewMethod that writes a SetBindingPhase method
// for the supplied Object to the supplied file.
func NewSetBindingPhase(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetBindingPhase of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetBindingPhase").Params(jen.Id("p").Qual(runtime, "BindingPhase")).Block(
			jen.Id(receiver).Dot(fields.NameStatus).Dot("SetBindingPhase").Call(jen.Id("p")),
		)
	}
}

// NewGetBindingPhase returns a NewMethod that writes a GetBindingPhase method
// for the supplied Object to the supplied file.
func NewGetBindingPhase(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetBindingPhase of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetBindingPhase").Params().Qual(runtime, "BindingPhase").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameStatus).Dot("GetBindingPhase").Call()),
		)
	}
}

// NewSetClaimReference returns a NewMethod that writes a SetClaimReference
// method for the supplied Object to the supplied file.
func NewSetClaimReference(receiver, core string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetClaimReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetClaimReference").Params(jen.Id("r").Op("*").Qual(core, "ObjectReference")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("ClaimReference").Op("=").Id("r"),
		)
	}
}

// NewGetClaimReference returns a NewMethod that writes a GetClaimReference
// method for the supplied Object to the supplied file.
func NewGetClaimReference(receiver, core string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetClaimReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetClaimReference").Params().Op("*").Qual(core, "ObjectReference").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("ClaimReference")),
		)
	}
}

// NewSetResourceReference returns a NewMethod that writes a
// SetResourceReference method for the supplied Object to the supplied file.
func NewSetResourceReference(receiver, core string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetResourceReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetResourceReference").Params(jen.Id("r").Op("*").Qual(core, "ObjectReference")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("ResourceReference").Op("=").Id("r"),
		)
	}
}

// NewGetResourceReference returns a NewMethod that writes a
// GetResourceReference method for the supplied Object to the supplied file.
func NewGetResourceReference(receiver, core string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetResourceReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetResourceReference").Params().Op("*").Qual(core, "ObjectReference").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("ResourceReference")),
		)
	}
}

// NewSetClassSelector returns a NewMethod that writes a SetClassSelector
// method for the supplied Object to the supplied file.
func NewSetClassSelector(receiver, meta string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetClassSelector of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetClassSelector").Params(jen.Id("s").Op("*").Qual(meta, "LabelSelector")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("ClassSelector").Op("=").Id("s"),
		)
	}
}

// NewGetClassSelector returns a NewMethod that writes a GetClassSelector
// method for the supplied Object to the supplied file.
func NewGetClassSelector(receiver, meta string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetClassSelector of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetClassSelector").Params().Op("*").Qual(meta, "LabelSelector").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("ClassSelector")),
		)
	}
}

// NewSetClassReference returns a NewMethod that writes a SetClassReference
// method for the supplied Object to the supplied file.
func NewSetClassReference(receiver, core string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetClassReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetClassReference").Params(jen.Id("r").Op("*").Qual(core, "ObjectReference")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("ClassReference").Op("=").Id("r"),
		)
	}
}

// NewGetClassReference returns a NewMethod that writes a GetClassReference
// method for the supplied Object to the supplied file.
func NewGetClassReference(receiver, core string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetClassReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetClassReference").Params().Op("*").Qual(core, "ObjectReference").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("ClassReference")),
		)
	}
}

// NewSetWriteConnectionSecretToReference returns a NewMethod that writes a
// SetWriteConnectionSecretToReference method for the supplied Object to the
// supplied file.
func NewSetWriteConnectionSecretToReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetWriteConnectionSecretToReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetWriteConnectionSecretToReference").Params(jen.Id("r").Op("*").Qual(runtime, "SecretReference")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("WriteConnectionSecretToReference").Op("=").Id("r"),
		)
	}
}

// NewGetWriteConnectionSecretToReference returns a NewMethod that writes a
// GetWriteConnectionSecretToReference method for the supplied Object to the
// supplied file.
func NewGetWriteConnectionSecretToReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetWriteConnectionSecretToReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetWriteConnectionSecretToReference").Params().Op("*").Qual(runtime, "SecretReference").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("WriteConnectionSecretToReference")),
		)
	}
}

// NewLocalSetWriteConnectionSecretToReference returns a NewMethod that writes a
// SetWriteConnectionSecretToReference method for the supplied Object to the
// supplied file.
func NewLocalSetWriteConnectionSecretToReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetWriteConnectionSecretToReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetWriteConnectionSecretToReference").Params(jen.Id("r").Op("*").Qual(runtime, "LocalSecretReference")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("WriteConnectionSecretToReference").Op("=").Id("r"),
		)
	}
}

// NewLocalGetWriteConnectionSecretToReference returns a NewMethod that writes a
// GetWriteConnectionSecretToReference method for the supplied Object to the
// supplied file.
func NewLocalGetWriteConnectionSecretToReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetWriteConnectionSecretToReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetWriteConnectionSecretToReference").Params().Op("*").Qual(runtime, "LocalSecretReference").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("WriteConnectionSecretToReference")),
		)
	}
}

// NewSetReclaimPolicy returns a NewMethod that writes a SetReclaimPolicy method
// for the supplied Object to the supplied file. The ReclaimPolicy is set in the
// supplied field - typically Spec or SpecTemplate.
func NewSetReclaimPolicy(receiver, core, field string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetReclaimPolicy of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetReclaimPolicy").Params(jen.Id("r").Qual(core, "ReclaimPolicy")).Block(
			jen.Id(receiver).Dot(field).Dot("ReclaimPolicy").Op("=").Id("r"),
		)
	}
}

// NewGetReclaimPolicy returns a NewMethod that writes a GetReclaimPolicy method
// for the supplied Object to the supplied file. The ReclaimPolicy is returned
// from the supplied field - typically Spec or SpecTemplate.
func NewGetReclaimPolicy(receiver, runtime, field string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetReclaimPolicy of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetReclaimPolicy").Params().Qual(runtime, "ReclaimPolicy").Block(
			jen.Return(jen.Id(receiver).Dot(field).Dot("ReclaimPolicy")),
		)
	}
}

// NewGetCredentialsSecretReference returns a NewMethod that writes a
// GetCredentialsSecretReference method for the supplied Object to the supplied file.
func NewGetCredentialsSecretReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetCredentialsSecretReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetCredentialsSecretReference").Params().Op("*").Qual(runtime, "SecretKeySelector").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("CredentialsSecretRef")),
		)
	}
}

// NewSetCredentialsSecretReference returns a NewMethod that writes a
// SetCredentialsSecretReference method for the supplied Object to the supplied file.
func NewSetCredentialsSecretReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetCredentialsSecretReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetCredentialsSecretReference").Params(jen.Id("r").Op("*").Qual(runtime, "SecretKeySelector")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("CredentialsSecretRef").Op("=").Id("r"),
		)
	}
}
