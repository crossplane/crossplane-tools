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

	"github.com/crossplane/crossplane-tools/internal/comments"
	"github.com/dave/jennifer/jen"
)

const (
	ReferenceTypeMarker               = "crossplane:generate:reference:type"
	ReferenceExtractorMarker          = "crossplane:generate:reference:extractor"
	ReferenceReferenceFieldNameMarker = "crossplane:generate:reference:refFieldName"
	ReferenceSelectorFieldNameMarker  = "crossplane:generate:reference:selectorFieldName"
)

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

func NewReferenceProcessor(defaultExtractor *jen.Statement) *ReferenceProcessor {
	return &ReferenceProcessor{
		DefaultExtractor: defaultExtractor,
	}
}

type ReferenceProcessor struct {
	DefaultExtractor *jen.Statement

	refs []Reference
}

func (rp *ReferenceProcessor) Process(_ *types.Named, f *types.Var, _ string, comment string, formerFields []string) error {
	markers := comments.ParseMarkers(comment)
	refTypeValues := markers[ReferenceTypeMarker]
	if len(refTypeValues) == 0 {
		return nil
	}
	refType := refTypeValues[0]
	isPointer := false
	isList := false
	// We don't support *[]string.
	switch t := f.Type().(type) {
	// *string
	case *types.Pointer:
		isPointer = true
	// []string.
	case *types.Slice:
		isList = true
		switch t.Elem().(type) {
		// []*string
		case *types.Pointer:
			isPointer = true
		}
	}

	extractorValues := markers[ReferenceExtractorMarker]
	extractorPath := rp.DefaultExtractor
	if len(extractorValues) != 0 {
		extractorPath = getTypeCodeFromPath(extractorValues[0])
	}
	fieldPath := strings.Join(append(formerFields, f.Name()), ".")
	rp.refs = append(rp.refs, Reference{
		RemoteType:          getTypeCodeFromPath(refType),
		RemoteListType:      getTypeCodeFromPath(refType, "List"),
		Extractor:           extractorPath,
		GoValueFieldPath:    fieldPath,
		GoRefFieldPath:      getRefFieldName(markers, fieldPath, isList),
		GoSelectorFieldPath: getSelectorFieldName(markers, fieldPath),
		IsPointer:           isPointer,
		IsList:              isList,
	})
	return nil
}

func (rp *ReferenceProcessor) GetReferences() []Reference {
	return rp.refs
}

func getRefFieldName(markers comments.Markers, valueFieldPath string, isList bool) string {
	if vals, ok := markers[ReferenceReferenceFieldNameMarker]; ok {
		f := strings.Split(valueFieldPath, ".")
		return strings.Join(f[:len(f)-1], ".") + "." + vals[0]
	}
	if isList {
		return valueFieldPath + "Refs"
	}
	return valueFieldPath + "Ref"
}

func getSelectorFieldName(markers comments.Markers, valueFieldPath string) string {
	if vals, ok := markers[ReferenceSelectorFieldNameMarker]; ok {
		f := strings.Split(valueFieldPath, ".")
		return strings.Join(f[:len(f)-1], ".") + "." + vals[0]
	}
	return valueFieldPath + "Selector"
}

func getTypeCodeFromPath(path string, nameSuffix ...string) *jen.Statement {
	words := strings.Split(path, ".")
	if len(words) == 1 {
		return jen.Op("&").Id(path + strings.Join(nameSuffix, "")).Values()
	}
	name := words[len(words)-1] + strings.Join(nameSuffix, "")
	pkg := strings.TrimSuffix(path, "."+words[len(words)-1])
	return jen.Op("&").Qual(pkg, name).Values()
}
