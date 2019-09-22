package generate

import (
	"go/parser"
	"go/token"
	"go/types"
	"io"

	"github.com/dave/jennifer/jen"
	"github.com/negz/angryjet/internal/match"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

type WriteOptions struct {
	Matches       match.Object
	ImportAliases map[string]string
	Headers       []string
}

type WriteOption func(o *WriteOptions)

func WithHeaders(h ...string) WriteOption {
	return func(o *WriteOptions) {
		o.Headers = append(o.Headers, h...)
	}
}

func WithMatcher(m match.Object) WriteOption {
	return func(o *WriteOptions) {
		o.Matches = m
	}
}

func WithImportAliases(ia map[string]string) WriteOption {
	return func(o *WriteOptions) {
		o.ImportAliases = ia
	}
}

func WriteMethods(p *packages.Package, m Methods, w io.Writer, wo ...WriteOption) error {
	opts := &WriteOptions{Matches: func(o types.Object) bool { return true }}
	for _, fn := range wo {
		fn(opts)
	}

	f := jen.NewFile(p.Name)
	for path, alias := range opts.ImportAliases {
		f.ImportAlias(path, alias)
	}
	for _, hc := range opts.Headers {
		f.HeaderComment(hc)
	}

	for _, n := range p.Types.Scope().Names() {
		o := p.Types.Scope().Lookup(n)
		if !opts.Matches(o) {
			continue
		}

		m.Generate(f, o)
	}

	return errors.Wrap(f.Render(w), "cannot render file")
}

type NewMethod func(f *jen.File, o types.Object)

type Methods map[string]NewMethod

func (m Methods) Generate(f *jen.File, o types.Object) {
	for name, gen := range m {
		if !match.HasMethod(name)(o) {
			gen(f, o)
		}
	}
}

func ProducedNothing(data []byte) bool {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "f.go", data, 0)
	if err != nil {
		return true
	}

	return len(f.Decls)+len(f.Scope.Objects) == 0
}
