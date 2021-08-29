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

	"github.com/pkg/errors"

	"github.com/crossplane/crossplane-tools/internal/comments"

	"github.com/dave/jennifer/jen"
)

// NewResolveReferences returns a NewMethod that writes a SetProviderConfigReference
// method for the supplied Object to the supplied file.
func NewResolveReferences(comm comments.Comments, refTypeMarker, refExtractorMarker, receiver, clientPath, referencePath string) New {
	return func(f *jen.File, o types.Object) {
		n, ok := o.Type().(*types.Named)
		if !ok {
			return
		}
		defaultExtractor := jen.Qual(referencePath, "ExternalName").Call()
		rs := NewReferenceSearcher(comm, defaultExtractor, refTypeMarker, refExtractorMarker)
		refs, err := rs.Search(n)
		if err != nil {
			panic(errors.Wrapf(err, "cannot search for references of %s", n.Obj().Name()))
		}
		if len(refs) == 0 {
			return
		}
		resolverCalls := make(jen.Statement, len(refs))
		for i, ref := range refs {
			currentValuePath := jen.Id(receiver)
			for _, fieldName := range strings.Split(ref.GoValueFieldPath, ".") {
				currentValuePath = currentValuePath.Dot(fieldName)
			}
			setResolvedValue := currentValuePath.Clone().Op("=").Id("resp").Dot("ResolvedValue")
			if ref.IsPointer {
				setResolvedValue = currentValuePath.Clone().Op("=").Qual(referencePath, "ToPtrValue").Call(jen.Id("resp").Dot("ResolvedValue"))
				currentValuePath = jen.Qual(referencePath, "FromPtrValue").Call(currentValuePath)
			}
			referenceFieldPath := jen.Id(receiver)
			for _, fieldName := range strings.Split(ref.GoRefFieldPath, ".") {
				referenceFieldPath = referenceFieldPath.Dot(fieldName)
			}
			selectorFieldPath := jen.Id(receiver)
			for _, fieldName := range strings.Split(ref.GoSelectorFieldPath, ".") {
				selectorFieldPath = selectorFieldPath.Dot(fieldName)
			}
			code := &jen.Statement{
				jen.List(jen.Id("resp"), jen.Err()).Op("=").Id("r").Dot("Resolve").Call(
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
				referenceFieldPath.Clone().Op("=").Id("resp").Dot("ResolvedReference"),
				jen.Line(),
			}
			resolverCalls[i] = code
		}

		f.Commentf("ResolveReferences of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("ResolveReferences").
			Params(
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id("c").Qual(clientPath, "Reader"),
			).Error().Block(
			jen.Id("r").Op(":=").Qual(referencePath, "NewAPIResolver").Call(jen.Id("c"), jen.Id(receiver)),
			jen.Line(),
			jen.Var().Id("resp").Qual(referencePath, "ResolutionResponse"),
			jen.Var().Err().Error(),
			jen.Line(),
			&resolverCalls,
			jen.Line(),
			jen.Return(jen.Nil()),
		)
	}
}

// Target type string

type Reference struct {
	RemoteType *jen.Statement
	Extractor  *jen.Statement

	RemoteListType      *jen.Statement
	GoValueFieldPath    string
	GoRefFieldPath      string
	GoSelectorFieldPath string
	IsList              bool
	IsPointer           bool
}

func NewReferenceSearcher(comm comments.Comments, defaultExtractor *jen.Statement, refTypeMarker, refExtractorMarker string) *ReferenceSearcher {
	return &ReferenceSearcher{
		Comments:                 comm,
		DefaultExtractor:         defaultExtractor,
		ReferenceTypeMarker:      refTypeMarker,
		ReferenceExtractorMarker: refExtractorMarker,
	}
}

type ReferenceSearcher struct {
	Comments                 comments.Comments
	ReferenceTypeMarker      string
	ReferenceExtractorMarker string
	DefaultExtractor         *jen.Statement

	refs []Reference
}

func (rs *ReferenceSearcher) Search(n *types.Named) ([]Reference, error) {
	return rs.refs, errors.Wrap(rs.search(n), "search for references failed")
}

func (rs *ReferenceSearcher) search(n *types.Named, fields ...string) error {
	s, ok := n.Underlying().(*types.Struct)
	if !ok {
		return nil
	}

	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)
		isPointer := false
		isList := false
		switch ft := field.Type().(type) {
		// Type
		case *types.Named:
			if err := rs.search(ft, append(fields, field.Name())...); err != nil {
				return errors.Wrapf(err, "cannot search for references in %s", ft.Obj().Name())
			}
		// *Type
		case *types.Pointer:
			isPointer = true
			switch elemType := ft.Elem().(type) {
			case *types.Named:
				if err := rs.search(elemType, append(fields, "*"+field.Name())...); err != nil {
					return errors.Wrapf(err, "cannot search for references in %s", elemType.Obj().Name())
				}
			}
		case *types.Slice:
			isList = true
			switch elemType := ft.Elem().(type) {
			// []Type
			case *types.Named:
				if err := rs.search(elemType, append(fields, "[]"+field.Name())...); err != nil {
					return errors.Wrapf(err, "cannot search for references in %s", elemType.Obj().Name())
				}
			// []*Type
			case *types.Pointer:
				switch elemElemType := elemType.Elem().(type) {
				case *types.Named:
					if err := rs.search(elemElemType, append(fields, "[]*"+field.Name())...); err != nil {
						return errors.Wrapf(err, "cannot search for references in %s", elemElemType.Obj().Name())
					}
				}
			}
		}
		markers := comments.ParseMarkers(rs.Comments.For(field))
		refTypeValues := markers[rs.ReferenceTypeMarker]
		if len(refTypeValues) == 0 {
			continue
		}
		refType := refTypeValues[0]

		extractorValues := markers[rs.ReferenceExtractorMarker]
		extractorPath := rs.DefaultExtractor
		if len(extractorValues) != 0 {
			extractorPath = getTypeCodeFromPath(extractorValues[0])
		}
		fieldPath := strings.Join(append(fields, field.Name()), ".")
		rs.refs = append(rs.refs, Reference{
			RemoteType:          getTypeCodeFromPath(refType),
			RemoteListType:      getTypeCodeFromPath(refType, "List"),
			Extractor:           extractorPath,
			GoValueFieldPath:    fieldPath,
			GoRefFieldPath:      fieldPath + "Ref",
			GoSelectorFieldPath: fieldPath + "Selector",
			IsPointer:           isPointer,
			IsList:              isList,
		})
	}
	return nil
}

func getTypeCodeFromPath(path string, nameSuffix ...string) *jen.Statement {
	words := strings.Split(path, ".")
	if len(words) == 1 {
		return jen.Op("&").Id(path + strings.Join(nameSuffix, "")).Values()
	}
	name := words[len(words)-1] + strings.Join(nameSuffix, "")
	pkg := strings.TrimSuffix(path, "."+name)
	return jen.Op("&").Qual(pkg, name).Values()
}
