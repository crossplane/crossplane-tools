package managed

import (
	"bytes"
	"go/types"
	"io/ioutil"
	"path/filepath"

	"github.com/dave/jennifer/jen"
	"github.com/negz/angryjet/internal/comments"
	"github.com/negz/angryjet/internal/fields"
	"github.com/negz/angryjet/internal/generate"
	"github.com/negz/angryjet/internal/match"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

const (
	RuntimeImport = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
	RuntimeAlias  = "runtimev1alpha1"
	Receiver      = "mg"
	FileSuffix    = ".managed.go"
	DisableMarker = "crossplane:generate:methods"
)

var Methods = generate.Methods{
	"SetBindingPhase": NewSetBindingPhase(Receiver, RuntimeImport),
	"GetBindingPhase": NewGetBindingPhase(Receiver, RuntimeImport),
}

func WriteMethods(base, prefix string, p *packages.Package) error {
	b := &bytes.Buffer{}
	err := generate.WriteMethods(p, Methods, b,
		generate.WithImportAliases(map[string]string{RuntimeImport: RuntimeAlias}),
		generate.WithMatcher(match.AllOf(Managed(), match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false"))),
	)
	if err != nil {
		return errors.Wrap(err, "cannot generate managed resource methods")
	}

	if generate.ProducedNothing(b.Bytes()) {
		return nil
	}

	f := filepath.Join(base, p.PkgPath, prefix+FileSuffix)
	return errors.Wrap(ioutil.WriteFile(f, b.Bytes(), 0644), "cannot write managed resource methods")
}

func Managed() match.Object {
	return func(o types.Object) bool {
		return fields.Has(o,
			fields.IsTypeMeta(),
			fields.IsObjectMeta(),
			fields.IsSpec().And(fields.HasFieldThat(fields.IsEmbedded(fields.IsResourceSpec()))),
			fields.IsStatus().And(fields.HasFieldThat(fields.IsEmbedded(fields.IsResourceStatus()))),
		)
	}
}

func NewSetBindingPhase(receiver, runtime string) generate.NewMethod {
	return func(f *jen.File, o types.Object) {
		f.Commentf("SetBindingPhase of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("SetBindingPhase").Params(jen.Id("p").Qual(runtime, "BindingPhase")).Block(
			jen.Id(receiver).Dot(fields.NameStatus).Dot("SetBindingPhase").Call(jen.Id("p")),
		)
	}
}

func NewGetBindingPhase(receiver, runtime string) generate.NewMethod {
	return func(f *jen.File, o types.Object) {
		f.Commentf("GetBindingPhase of this %s.", o.Name())
		f.Func().Params(jen.Id(receiver).Op("*").Id(o.Name())).Id("GetBindingPhase").Params().Qual(runtime, "BindingPhase").Block(
			jen.Return(jen.Id(receiver).Dot(fields.NameStatus).Dot("GetBindingPhase").Call()),
		)
	}
}
