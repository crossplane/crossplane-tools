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

package types

import (
	"go/types"

	"github.com/pkg/errors"

	"github.com/crossplane/crossplane-tools/internal/comments"
)

// NamedProcessorChain runs multiple NamedProcessors in order.
type NamedProcessorChain []NamedProcessor

// Process run Process method of all NamedProcessors.
func (tpc NamedProcessorChain) Process(n *types.Named, comment string) error {
	for i, tp := range tpc {
		if err := tp.Process(n, comment); err != nil {
			return errors.Errorf("type processor at index %d failed", i)
		}
	}
	return nil
}

// NamedProcessor takes a named struct with its comments and processes it. For
// example, you might want to take an action depending on its comments or make
// an in-place change on the *types.Named object if it satisfies certain criteria.
type NamedProcessor interface {
	Process(n *types.Named, comment string) error
}

// FieldProcessorChain runs multiple FieldProcessor in order.
type FieldProcessorChain []FieldProcessor

// Process run Process method of all FieldProcessor.
func (fpc FieldProcessorChain) Process(n *types.Named, f *types.Var, tag, comment string, parentFields ...string) error {
	for i, fp := range fpc {
		if err := fp.Process(n, f, tag, comment, parentFields...); err != nil {
			return errors.Errorf("field processor at index %d failed", i)
		}
	}
	return nil
}

// FieldProcessor takes all information related to a field and processes it.
type FieldProcessor interface {
	Process(n *types.Named, f *types.Var, tag string, comment string, parentFields ...string) error
}

// ProcessorConfig lets you configure what processors will be run in given traversal.
type ProcessorConfig struct {
	Named NamedProcessor
	Field FieldProcessor
}

// NewTraverser returns a new Traverser.
func NewTraverser(c comments.Comments) *Traverser {
	return &Traverser{
		comments: c,
	}
}

// Traverser goes through all fields of given type recursively. It runs the field
// processor for every field and named processor for every type it encounters
// during its depth-first traversal.
type Traverser struct {
	comments comments.Comments
}

// NOTE(muvaf): We return an error but currently there isn't really anything
// constructing an error. But we keep that for future type and field processors.

// Traverse traverser given type recursively and runs given processors.
func (t *Traverser) Traverse(n *types.Named, cfg *ProcessorConfig, parentFields ...string) error { // nolint:gocyclo
	// NOTE(muvaf): gocyclo is disabled due to repeated type checks.
	if err := cfg.Named.Process(n, t.comments.For(n.Obj())); err != nil {
		return errors.Wrapf(err, "type processors failed to run for type %s", n.Obj().Name())
	}
	st, ok := n.Underlying().(*types.Struct)
	if !ok {
		return nil
	}
	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		tag := st.Tag(i)
		if err := cfg.Field.Process(n, field, tag, t.comments.For(field), parentFields...); err != nil {
			return errors.Wrapf(err, "field processors failed to run for field %s of type %s", field.Name(), n.Obj().Name())
		}
		switch ft := field.Type().(type) {
		case *types.Named:
			if err := t.Traverse(ft, cfg, append(parentFields, field.Name())...); err != nil {
				return errors.Wrapf(err, "failed to traverse type of field %s", field.Name())
			}
		case *types.Pointer:
			if elemType, ok := ft.Elem().(*types.Named); ok {
				if err := t.Traverse(elemType, cfg, append(parentFields, "*"+field.Name())...); err != nil {
					return errors.Wrapf(err, "failed to traverse type of field %s", field.Name())
				}
			}
		case *types.Slice:
			switch elemType := ft.Elem().(type) {
			case *types.Named:
				if err := t.Traverse(elemType, cfg, append(parentFields, "[]"+field.Name())...); err != nil {
					return errors.Wrapf(err, "failed to traverse type of field %s", field.Name())
				}
			case *types.Pointer:
				if elemElemType, ok := elemType.Elem().(*types.Named); ok {
					if err := t.Traverse(elemElemType, cfg, append(parentFields, "[]"+"*"+field.Name())...); err != nil {
						return errors.Wrapf(err, "failed to traverse type of field %s", field.Name())
					}
				}
			}
		}
	}
	return nil
}
