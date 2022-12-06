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
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"

	"github.com/crossplane/crossplane-tools/internal/comments"
)

// Comment markers used by ReferenceProcessor
const (
	ReferenceTypeMarker               = "crossplane:generate:reference:type"
	ReferenceExtractorMarker          = "crossplane:generate:reference:extractor"
	ReferenceReferenceFieldNameMarker = "crossplane:generate:reference:refFieldName"
	ReferenceSelectorFieldNameMarker  = "crossplane:generate:reference:selectorFieldName"
)

var (
	regexFunctionCall = regexp.MustCompile(`((.+)\.)?([^.]+\(.*\))`)
)

// Reference is the internal representation that has enough information to let
// us generate the resolver.
type Reference struct {
	// RemoteType represents the type whose reference we're holding.
	RemoteType *jen.Statement

	// Extractor is the function call of the function that will take referenced
	// instance and return a string or []string to be set as value.
	Extractor *jen.Statement

	// RemoteListType is the list type of the type whose reference we're holding.
	RemoteListType *jen.Statement

	// GoValueFieldPath is the list of fields that needs to be traveled to access
	// the current value field. It may include prefixes like [] for array fields,
	// * for pointer fields or []* for array of pointer fields.
	GoValueFieldPath []string

	// GoRefFieldName is the name of the field whose type is *xpv1.Reference or
	// []xpv1.Reference.
	GoRefFieldName string

	// GoSelectorFieldName is the name of the field whose type is *xpv1.Selector
	GoSelectorFieldName string

	// IsSlice tells whether the current value type is a slice kind.
	IsSlice bool

	// IsPointer tells whether the current value type is a pointer kind.
	IsPointer bool

	// IsFloatPointer tells whether the current value pointer is of type float64
	IsFloatPointer bool
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
func (rp *ReferenceProcessor) Process(_ *types.Named, f *types.Var, _, comment string, parentFields ...string) error {
	markers := comments.ParseMarkers(comment)
	refTypeValues := markers[ReferenceTypeMarker]
	if len(refTypeValues) == 0 {
		return nil
	}
	refType := refTypeValues[0]
	isPointer := false
	isList := false
	isFloatPointer := false
	refFieldName := f.Name() + "Ref"

	// We don't support *[]string.
	switch t := f.Type().(type) {
	// *string|*float64
	case *types.Pointer:
		isPointer = true
	// []string|[]float64.
	case *types.Slice:
		isList = true
		refFieldName = f.Name() + "Refs"
		// []*string|[]*float64
		if _, ok := t.Elem().(*types.Pointer); ok {
			isPointer = true
		}
	}

	if strings.HasSuffix(f.Type().String(), "*float64") {
		isFloatPointer = true
	}

	extractorPath := rp.DefaultExtractor
	if values, ok := markers[ReferenceExtractorMarker]; ok {
		var err error
		extractorPath, err = getFuncCodeFromPath(values[0])
		if err != nil {
			return errors.Wrapf(err, "cannot get extractor function")
		}
	}

	if values, ok := markers[ReferenceReferenceFieldNameMarker]; ok {
		refFieldName = values[0]
	}

	selectorFieldName := f.Name() + "Selector"
	if values, ok := markers[ReferenceSelectorFieldNameMarker]; ok {
		selectorFieldName = values[0]
	}
	path := append([]string{rp.Receiver}, parentFields...)
	rp.refs = append(rp.refs, Reference{
		RemoteType:          getTypeCodeFromPath(refType),
		RemoteListType:      getTypeCodeFromPath(refType, "List"),
		Extractor:           extractorPath,
		GoValueFieldPath:    append(path, f.Name()),
		GoRefFieldName:      refFieldName,
		GoSelectorFieldName: selectorFieldName,
		IsPointer:           isPointer,
		IsSlice:             isList,
		IsFloatPointer:      isFloatPointer,
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

func getFuncCodeFromPath(path string) (*jen.Statement, error) {
	parts := regexFunctionCall.FindStringSubmatch(path)
	// we have a total of four groups in the regular expression so if
	// we do not have four parts, then we cannot handle the reference expression
	// Examples paths are:
	// github.com/upbound/upjet/pkg/resource.ExtractParamPath("a.b.c",true)
	// ExtractParamPath("a.b.c",true)
	// ExtractParamPath("a", false)
	// ExtractParamPath()
	if len(parts) != 4 {
		return nil, errors.Errorf("path %q is not a valid function code", path)
	}
	if len(parts[1]) == 0 {
		return jen.Id(parts[3]), nil
	}
	return jen.Qual(parts[2], parts[3]), nil
}
