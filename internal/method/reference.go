/*
Copyright 2021 The Crossplane Authors.

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
	"go/types"

	"github.com/dave/jennifer/jen"
)

// NewResolveReferences returns a NewMethod that writes a SetProviderConfigReference
// method for the supplied Object to the supplied file.
func NewResolveReferences(receiver, clientPath, referencePath string) New {
	return func(f *jen.File, o types.Object) {
		f.Commentf("ResolveReferences of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("ResolveReferences").
			Params(
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id("c").Qual(clientPath, "Reader"),
			).Block(
			jen.Id("r").Op(":=").Qual(referencePath, "NewAPIResolver").Call(jen.Id("c"), jen.Id(receiver)),
			jen.Line(),
			jen.List(jen.Id("resp"), jen.Err()).Op(":=").Id("r").Dot("Resolve").Call(
				jen.Id("ctx"),
				jen.Qual(referencePath, "ResolutionRequest").Values(jen.Dict{
					jen.Id("CurrentValue"): jen.Id(receiver).Dot("Spec").Dot("ForProvider").Dot("ApiId"),
					jen.Id("Reference"):    jen.Id(receiver).Dot("Spec").Dot("ForProvider").Dot("ApiIdRef"),
					jen.Id("Selector"):     jen.Id(receiver).Dot("Spec").Dot("ForProvider").Dot("ApiIdSelector"),
					jen.Id("To"): jen.Qual(referencePath, "To").Values(jen.Dict{
						jen.Id("Managed"): jen.Op("&").Id("Apigatewayv2Api").Values(),
						jen.Id("List"):    jen.Op("&").Id("Apigatewayv2ApiList").Values(),
					}),
					jen.Id("Extract"): jen.Qual(referencePath, "ExternalName").Call(),
				},
				),
			),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Qual("github.com/pkg/errors", "Wrap").Call(jen.Err(), jen.Lit("Spec.ForProvider.ApiId"))),
			),
			jen.Line(),
			jen.Id(receiver).Dot("Spec").Dot("ForProvider").Dot("ApiId").Op("=").Id("resp").Dot("ResolvedValue"),
			jen.Id(receiver).Dot("Spec").Dot("ForProvider").Dot("ApiIdRef").Op("=").Id("resp").Dot("ResolvedReference"),
			jen.Line(),
			jen.Return(jen.Nil()),
		)
	}
}
