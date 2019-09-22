// Package managed writes the method set required to satisfy
// crossplane-runtime's resource.Managed interface.
package managed

import (
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"

	"github.com/negz/angryjet/internal/comments"
	"github.com/negz/angryjet/internal/generate"
	"github.com/negz/angryjet/internal/generators/methods"
	"github.com/negz/angryjet/internal/match"
)

// Imports used in generated code.
const (
	CoreAlias  = "corev1"
	CoreImport = "k8s.io/api/core/v1"

	RuntimeAlias  = "runtimev1alpha1"
	RuntimeImport = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
)

const (
	// Receiver name for managed resource generated methods.
	Receiver = "mg"

	// FileName for managed resource generated methods.
	FileName = "managed.go"

	// DisableMarker used to disable generation of managed resource methods for
	// a type that otherwise appears to be a managed resource that is missing a
	// subnet of its methods.
	DisableMarker = "crossplane:generate:methods"
)

// Methods to generate for managed resources that are missing them.
var Methods = generate.MethodSet{
	"SetConditions":                       methods.NewSetConditions(Receiver, RuntimeImport),
	"SetBindingPhase":                     methods.NewSetBindingPhase(Receiver, RuntimeImport),
	"GetBindingPhase":                     methods.NewGetBindingPhase(Receiver, RuntimeImport),
	"SetClaimReference":                   methods.NewSetClaimReference(Receiver, CoreImport),
	"GetClaimReference":                   methods.NewGetClaimReference(Receiver, CoreImport),
	"SetNonPortableClassReference":        methods.NewSetNonPortableClassReference(Receiver, CoreImport),
	"GetNonPortableClassReference":        methods.NewGetNonPortableClassReference(Receiver, CoreImport),
	"SetWriteConnectionSecretToReference": methods.NewSetWriteConnectionSecretToReference(Receiver, CoreImport),
	"GetWriteConnectionSecretToReference": methods.NewGetWriteConnectionSecretToReference(Receiver, CoreImport),
	"SetReclaimPolicy":                    methods.NewSetReclaimPolicy(Receiver, RuntimeImport),
	"GetReclaimPolicy":                    methods.NewGetReclaimPolicy(Receiver, RuntimeImport),
}

// WriteMethods for all managed resources in the supplied package. Methods will
// be written to a file under the package's path relative to the supplied base
// path. The file will be named FileName, prefixed with the supplied prefix. The
// supplied header will be written as a comment to all generated files.
func WriteMethods(base, prefix, header string, p *packages.Package) error {
	err := generate.WriteMethods(p, Methods, filepath.Join(base, p.PkgPath, prefix+FileName),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{
			CoreImport:    CoreAlias,
			RuntimeImport: RuntimeAlias,
		}),
		generate.WithMatcher(match.AllOf(
			match.Managed(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)
	return errors.Wrap(err, "cannot write managed resource methods")
}
