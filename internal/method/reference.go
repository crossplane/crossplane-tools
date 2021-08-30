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
	"strings"

	twpackages "github.com/muvaf/typewriter/pkg/packages"

	"github.com/dave/jennifer/jen"
	twtypes "github.com/muvaf/typewriter/pkg/types"
)

// NewResolveReferences returns a NewMethod that writes a SetProviderConfigReference
// method for the supplied Object to the supplied file.
func NewResolveReferences(cache *twpackages.Cache, receiver, clientPath, referencePath string) New {
	return func(f *jen.File, o types.Object) {
		n, ok := o.Type().(*types.Named)
		if !ok {
			return
		}
		defaultExtractor := jen.Qual(referencePath, "ExternalName").Call()
		refProcessor := NewReferenceProcessor(defaultExtractor)
		tr := twtypes.NewTraverser(cache,
			twtypes.WithFieldProcessor(refProcessor),
		)
		if err := tr.Traverse(n); err != nil {
			panic("cannot traverse the type tree")
		}
		refs := refProcessor.GetReferences()
		if len(refs) == 0 {
			return
		}
		hasMultiResolution := false
		hasSingleResolution := false
		resolverCalls := make(jen.Statement, len(refs))
		for i, ref := range refs {
			if ref.IsList {
				hasMultiResolution = true
			} else {
				hasSingleResolution = true
			}
			currentValuePath := jen.Id(receiver)
			for _, fieldName := range strings.Split(ref.GoValueFieldPath, ".") {
				currentValuePath = currentValuePath.Dot(fieldName)
			}
			referenceFieldPath := jen.Id(receiver)
			for _, fieldName := range strings.Split(ref.GoRefFieldPath, ".") {
				referenceFieldPath = referenceFieldPath.Dot(fieldName)
			}
			selectorFieldPath := jen.Id(receiver)
			for _, fieldName := range strings.Split(ref.GoSelectorFieldPath, ".") {
				selectorFieldPath = selectorFieldPath.Dot(fieldName)
			}
			var code *jen.Statement
			if ref.IsList {
				code = &jen.Statement{
					jen.List(jen.Id("mrsp"), jen.Err()).Op("=").Id("r").Dot("ResolveMultiple").Call(
						jen.Id("ctx"),
						jen.Qual(referencePath, "MultiResolutionRequest").Values(jen.Dict{
							jen.Id("CurrentValues"): currentValuePath,
							jen.Id("References"):    referenceFieldPath,
							jen.Id("Selector"):      selectorFieldPath,
							jen.Id("To"): jen.Qual(referencePath, "To").Values(jen.Dict{
								jen.Id("Managed"): ref.RemoteType,
								jen.Id("List"):    ref.RemoteListType,
							}),
							jen.Id("Extract"): ref.Extractor,
						},
						),
					),
					jen.Line(),
					jen.If(jen.Err().Op("!=").Nil()).Block(
						jen.Return(jen.Qual("github.com/pkg/errors", "Wrapf").Call(jen.Err(), jen.Lit(ref.GoValueFieldPath))),
					),
					jen.Line(),
					currentValuePath.Clone().Op("=").Id("mrsp").Dot("ResolvedValues"),
					jen.Line(),
					referenceFieldPath.Clone().Op("=").Id("mrsp").Dot("ResolvedReferences"),
					jen.Line(),
				}
			} else {
				setResolvedValue := currentValuePath.Clone().Op("=").Id("rsp").Dot("ResolvedValue")
				if ref.IsPointer {
					setResolvedValue = currentValuePath.Clone().Op("=").Qual(referencePath, "ToPtrValue").Call(jen.Id("rsp").Dot("ResolvedValue"))
					currentValuePath = jen.Qual(referencePath, "FromPtrValue").Call(currentValuePath)
				}
				code = &jen.Statement{
					jen.List(jen.Id("rsp"), jen.Err()).Op("=").Id("r").Dot("Resolve").Call(
						jen.Id("ctx"),
						jen.Qual(referencePath, "ResolutionRequest").Values(jen.Dict{
							jen.Id("CurrentValue"): currentValuePath,
							jen.Id("Reference"):    referenceFieldPath,
							jen.Id("Selector"):     selectorFieldPath,
							jen.Id("To"): jen.Qual(referencePath, "To").Values(jen.Dict{
								jen.Id("Managed"): ref.RemoteType,
								jen.Id("List"):    ref.RemoteListType,
							}),
							jen.Id("Extract"): ref.Extractor,
						},
						),
					),
					jen.Line(),
					jen.If(jen.Err().Op("!=").Nil()).Block(
						jen.Return(jen.Qual("github.com/pkg/errors", "Wrapf").Call(jen.Err(), jen.Lit(ref.GoValueFieldPath))),
					),
					jen.Line(),
					setResolvedValue,
					jen.Line(),
					referenceFieldPath.Clone().Op("=").Id("rsp").Dot("ResolvedReference"),
					jen.Line(),
				}
			}
			resolverCalls[i] = code
		}
		var initStatements jen.Statement
		if hasSingleResolution {
			initStatements = append(initStatements, jen.Var().Id("rsp").Qual(referencePath, "ResolutionResponse"), jen.Line())
		}
		if hasMultiResolution {
			initStatements = append(initStatements, jen.Var().Id("mrsp").Qual(referencePath, "MultiResolutionResponse"))
		}

		f.Commentf("ResolveReferences of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("ResolveReferences").
			Params(
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id("c").Qual(clientPath, "Reader"),
			).Error().Block(
			jen.Id("r").Op(":=").Qual(referencePath, "NewAPIResolver").Call(jen.Id("c"), jen.Id(receiver)),
			jen.Line(),
			&initStatements,
			jen.Var().Err().Error(),
			jen.Line(),
			&resolverCalls,
			jen.Line(),
			jen.Return(jen.Nil()),
		)
	}
}
