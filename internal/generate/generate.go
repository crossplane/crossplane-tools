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

// Package generate generates method sets for Go types.
package generate

import (
	"bytes"
	"go/parser"
	"go/token"
	"go/types"
	"os"

	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"

	"github.com/crossplane/crossplane-tools/internal/match"
	"github.com/crossplane/crossplane-tools/internal/method"
)

// HeaderGenerated is added to all files generated by angryjet.
// See https://github.com/golang/go/issues/13560#issuecomment-288457920.
const HeaderGenerated = "Code generated by angryjet. DO NOT EDIT."

type options struct {
	Matches       match.Object
	ImportAliases map[string]string
	Headers       []string
}

// A WriteOption configures method generation behaviour.
type WriteOption func(o *options)

// WithHeaders specifies strings to be written as comments to the generated
// file, above the package definition. Single line strings use // comments,
// while multiline strings use /**/ comments.
func WithHeaders(h ...string) WriteOption {
	return func(o *options) {
		o.Headers = append(o.Headers, h...)
	}
}

// WithMatcher specifies an Object matcher that is used to filter the Objects
// within the package down to the set that need the generated methods.
func WithMatcher(m match.Object) WriteOption {
	return func(o *options) {
		o.Matches = m
	}
}

// WithImportAliases configures a map of import paths to aliases that will be
// used when generating code. For example if a generated method requires
// "example.org/foo/bar" it may refer to that package as "foobar" by supplying
// map[string]string{"example.org/foo/bar": "foobar"}.
func WithImportAliases(ia map[string]string) WriteOption {
	return func(o *options) {
		o.ImportAliases = ia
	}
}

// WriteMethods writes the supplied methods for each object in the supplied
// package to the supplied file. Use WithMatcher to limit the objects for which
// methods will be written. Methods will not be generated if a method with the
// same name is already defined for the object outside of the supplied filename.
// Files will not be written if they would contain no methods.
func WriteMethods(p *packages.Package, ms method.Set, file string, wo ...WriteOption) error {
	opts := &options{Matches: func(_ types.Object) bool { return true }}
	for _, fn := range wo {
		fn(opts)
	}

	// NewFilePath creates a new File object by taking the full package path such as:
	// 'github.com/org/repo/apis/resource/v1alpha1'
	// File object created using the function ('NewFile') that takes only the package
	// name ('v1alpha1') is not sufficient in order to properly handle imports from
	// the same package and to communicate with the Jennifer tool correctly.
	// We need to create the File object by using NewFilePath (passing package path)
	// so that we can communicate correctly with the library (jennifer).
	f := jen.NewFilePath(p.PkgPath)
	for path, alias := range opts.ImportAliases {
		f.ImportAlias(path, alias)
	}
	for _, hc := range opts.Headers {
		if hc != "" {
			f.HeaderComment(hc)
		}
	}
	f.HeaderComment(HeaderGenerated)

	for _, n := range p.Types.Scope().Names() {
		o := p.Types.Scope().Lookup(n)
		if !opts.Matches(o) {
			continue
		}
		ms.Write(f, o, method.DefinedOutside(p.Fset, file))
	}

	b := &bytes.Buffer{}
	if err := f.Render(b); err != nil {
		return errors.Wrap(err, "cannot render Go file")
	}

	if ProducedNothing(b.Bytes()) {
		return nil
	}

	return errors.Wrap(os.WriteFile(file, b.Bytes(), 0o644), "cannot write Go file") //nolint:gosec // We're comfortable with this being world readable.
}

// ProducedNothing returns true if the supplied data is either not a valid Go
// source file, or a valid Go file that contains no top level objects or
// declarations.
func ProducedNothing(data []byte) bool {
	f, err := parser.ParseFile(token.NewFileSet(), "f.go", data, 0)
	if err != nil {
		return true
	}
	return len(f.Decls)+len(f.Scope.Objects) == 0
}
