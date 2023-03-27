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

// NewSetProviderReference returns a NewMethod that writes a SetProviderReference
// method for the supplied Object to the supplied file.
func NewSetProviderReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetProviderReference of this %s.\nDeprecated: Use SetProviderConfigReference.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetProviderReference").Params(jen.Id("r").Op("*").Qual(runtime, "Reference")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("ProviderReference").Op("=").Id("r"),
		)
	}
}

// NewGetProviderReference returns a NewMethod that writes a GetProviderReference
// method for the supplied Object to the supplied file.
func NewGetProviderReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetProviderReference of this %s.\nDeprecated: Use GetProviderConfigReference.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetProviderReference").Params().Op("*").Qual(runtime, "Reference").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("ProviderReference")),
		)
	}
}

// NewSetProviderConfigReference returns a NewMethod that writes a SetProviderConfigReference
// method for the supplied Object to the supplied file.
func NewSetProviderConfigReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetProviderConfigReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetProviderConfigReference").Params(jen.Id("r").Op("*").Qual(runtime, "Reference")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("ProviderConfigReference").Op("=").Id("r"),
		)
	}
}

// NewGetProviderConfigReference returns a NewMethod that writes a GetProviderConfigReference
// method for the supplied Object to the supplied file.
func NewGetProviderConfigReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetProviderConfigReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetProviderConfigReference").Params().Op("*").Qual(runtime, "Reference").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("ProviderConfigReference")),
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

// NewSetPublishConnectionDetailsTo returns a NewMethod that writes a
// NewSetPublishConnectionDetailsTo method for the supplied Object to the
// supplied file.
func NewSetPublishConnectionDetailsTo(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetPublishConnectionDetailsTo of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetPublishConnectionDetailsTo").Params(jen.Id("r").Op("*").Qual(runtime, "PublishConnectionDetailsTo")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("PublishConnectionDetailsTo").Op("=").Id("r"),
		)
	}
}

// NewGetPublishConnectionDetailsTo returns a NewMethod that writes a
// GetPublishConnectionDetailsTo method for the supplied Object to the
// supplied file.
func NewGetPublishConnectionDetailsTo(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetPublishConnectionDetailsTo of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetPublishConnectionDetailsTo").Params().Op("*").Qual(runtime, "PublishConnectionDetailsTo").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("PublishConnectionDetailsTo")),
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

// NewSetManagementPolicy returns a NewMethod that writes a SetManagementPolicy
// method for the supplied Object to the supplied file.
func NewSetManagementPolicy(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetManagementPolicy of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetManagementPolicy").Params(jen.Id("r").Qual(runtime, "ManagementPolicy")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("ManagementPolicy").Op("=").Id("r"),
		)
	}
}

// NewGetManagementPolicy returns a NewMethod that writes a GetManagementPolicy
// method for the supplied Object to the supplied file.
func NewGetManagementPolicy(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetManagementPolicy of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetManagementPolicy").Params().Qual(runtime, "ManagementPolicy").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("ManagementPolicy")),
		)
	}
}

// NewSetDeletionPolicy returns a NewMethod that writes a SetDeletionPolicy
// method for the supplied Object to the supplied file.
func NewSetDeletionPolicy(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetDeletionPolicy of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetDeletionPolicy").Params(jen.Id("r").Qual(runtime, "DeletionPolicy")).Block(
			jen.Id(receiver).Dot(fields.NameSpec).Dot("DeletionPolicy").Op("=").Id("r"),
		)
	}
}

// NewGetDeletionPolicy returns a NewMethod that writes a GetDeletionPolicy
// method for the supplied Object to the supplied file.
func NewGetDeletionPolicy(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetDeletionPolicy of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetDeletionPolicy").Params().Qual(runtime, "DeletionPolicy").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameSpec).Dot("DeletionPolicy")),
		)
	}
}

// NewSetUsers returns a NewMethod that writes a SetUsers method for the
// supplied Object to the supplied file.
func NewSetUsers(receiver string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetUsers of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetUsers").Params(jen.Id("i").Int64()).Block(
			jen.Id(receiver).Dot(fields.NameStatus).Dot("Users").Op("=").Id("i"),
		)
	}
}

// NewGetUsers returns a NewMethod that writes a GetUsers method for the
// supplied Object to the supplied file.
func NewGetUsers(receiver string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetUsers of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetUsers").Params().Int64().Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameStatus).Dot("Users")),
		)
	}
}

// NewManagedGetItems returns a New that writes a GetItems method for the
// supplied object to the supplied file.
func NewManagedGetItems(receiver, resource string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetItems of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetItems").Params().Index().Qual(resource, "Managed").Block(
			jen.Id("items").Op(":=").Make(jen.Index().Qual(resource, "Managed"), jen.Len(jen.Id(receiver).Dot("Items"))),
			jen.For(jen.Id("i").Op(":=").Range().Id(receiver).Dot("Items")).Block(
				jen.Id("items").Index(jen.Id("i")).Op("=").Op("&").Id(receiver).Dot("Items").Index(jen.Id("i")),
			),
			jen.Return(jen.Id("items")),
		)
	}
}

// NewSetRootProviderConfigReference returns a NewMethod that writes a
// SetProviderConfigReference method for the supplied Object to the supplied
// file. Note that unlike NewSetProviderConfigReference the generated method
// expects the ProviderConfigReference to be at the root of the struct, not
// under its Spec field.
func NewSetRootProviderConfigReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetProviderConfigReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetProviderConfigReference").Params(jen.Id("r").Qual(runtime, "Reference")).Block(
			jen.Id(receiver).Dot("ProviderConfigReference").Op("=").Id("r"),
		)
	}
}

// NewGetRootProviderConfigReference returns a NewMethod that writes a
// GetProviderConfigReference method for the supplied Object to the supplied
// file. file. Note that unlike NewGetProviderConfigReference the generated
// method expects the ProviderConfigReference to be at the root of the struct,
// not under its Spec field.
func NewGetRootProviderConfigReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetProviderConfigReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetProviderConfigReference").Params().Qual(runtime, "Reference").Block(
			jen.Return(jen.Id(receiver).Dot("ProviderConfigReference")),
		)
	}
}

// NewSetRootResourceReference returns a NewMethod that writes a
// SetRootResourceReference method for the supplied Object to the supplied file.
func NewSetRootResourceReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetResourceReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetResourceReference").Params(jen.Id("r").Qual(runtime, "TypedReference")).Block(
			jen.Id(receiver).Dot("ResourceReference").Op("=").Id("r"),
		)
	}
}

// NewGetRootResourceReference returns a NewMethod that writes a
// GetRootResourceReference method for the supplied Object to the supplied file.
func NewGetRootResourceReference(receiver, runtime string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetResourceReference of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetResourceReference").Params().Qual(runtime, "TypedReference").Block(
			jen.Return(jen.Id(receiver).Dot("ResourceReference")),
		)
	}
}

// NewProviderConfigUsageGetItems returns a New that writes a GetItems method for the
// supplied object to the supplied file.
func NewProviderConfigUsageGetItems(receiver, resource string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetItems of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetItems").Params().Index().Qual(resource, "ProviderConfigUsage").Block(
			jen.Id("items").Op(":=").Make(jen.Index().Qual(resource, "ProviderConfigUsage"), jen.Len(jen.Id(receiver).Dot("Items"))),
			jen.For(jen.Id("i").Op(":=").Range().Id(receiver).Dot("Items")).Block(
				jen.Id("items").Index(jen.Id("i")).Op("=").Op("&").Id(receiver).Dot("Items").Index(jen.Id("i")),
			),
			jen.Return(jen.Id("items")),
		)
	}
}
