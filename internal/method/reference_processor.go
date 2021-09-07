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

	"github.com/dave/jennifer/jen"

	"github.com/crossplane/crossplane-tools/internal/comments"
)

// Comment markers used by ReferenceProcessor
const (
	ReferenceTypeMarker               = "crossplane:generate:reference:type"
	ReferenceExtractorMarker          = "crossplane:generate:reference:extractor"
	ReferenceReferenceFieldNameMarker = "crossplane:generate:reference:refFieldName"
	ReferenceSelectorFieldNameMarker  = "crossplane:generate:reference:selectorFieldName"
)

// Reference is the internal representation that has enough information to let
// us generate the resolver.
type Reference struct {
	RemoteType *jen.Statement
	Extractor  *jen.Statement

	RemoteListType      *jen.Statement
	GoValueFieldPath    []string
	GoRefFieldName      string
	GoSelectorFieldName string
	IsList              bool
	IsPointer           bool
}

// ReferenceProcessorOption is used to configure ReferenceProcessor.
type ReferenceProcessorOption func(*ReferenceProcessor)

// WithDefaultExtractor returns an option that sets the extractor to given
// call.
func WithDefaultExtractor(ext *jen.Statement) ReferenceProcessorOption {
	return func(rp *ReferenceProcessor) {
		rp.DefaultExtractor = ext
	}
}

// NewReferenceProcessor returns a new *ReferenceProcessor .
func NewReferenceProcessor(receiver string, opts ...ReferenceProcessorOption) *ReferenceProcessor {
	rp := &ReferenceProcessor{
		Receiver: receiver,
	}
	for _, f := range opts {
		f(rp)
	}
	return rp
}

// ReferenceProcessor detects whether the field is marked as referencer and
// composes the internal representation of that reference.
type ReferenceProcessor struct {
	// DefaultExtractor is used when the extractor is not overridden.
	DefaultExtractor *jen.Statement

	// Receiver is prepended to all field paths.
	Receiver string

	refs []Reference
}

// Process stores the reference information of the given field, if any.
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
		// []*string
		if _, ok := t.Elem().(*types.Pointer); ok {
			isPointer = true
		}
	}

	extractorPath := rp.DefaultExtractor
	if values, ok := markers[ReferenceExtractorMarker]; ok {
		extractorPath = getFuncCodeFromPath(values[0])
	}

	refFieldName := f.Name() + "Ref"
	if isList {
		refFieldName = f.Name() + "Refs"
	}
	if values, ok := markers[ReferenceReferenceFieldNameMarker]; ok {
		refFieldName = values[0]
	}

	selectorFieldName := f.Name() + "Selector"
	if values, ok := markers[ReferenceSelectorFieldNameMarker]; ok {
		selectorFieldName = values[0]
	}
	path := append([]string{rp.Receiver}, formerFields...)
	rp.refs = append(rp.refs, Reference{
		RemoteType:          getTypeCodeFromPath(refType),
		RemoteListType:      getTypeCodeFromPath(refType, "List"),
		Extractor:           extractorPath,
		GoValueFieldPath:    append(path, f.Name()),
		GoRefFieldName:      refFieldName,
		GoSelectorFieldName: selectorFieldName,
		IsPointer:           isPointer,
		IsList:              isList,
	})
	return nil
}

// GetReferences returns all the references accumulated so far from processing.
func (rp *ReferenceProcessor) GetReferences() []Reference {
	return rp.refs
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

func getFuncCodeFromPath(path string) *jen.Statement {
	words := strings.Split(path, ".")
	if len(words) == 1 {
		return jen.Id(path)
	}
	name := words[len(words)-1]
	pkg := strings.TrimSuffix(path, "."+words[len(words)-1])
	return jen.Qual(pkg, name)
}
