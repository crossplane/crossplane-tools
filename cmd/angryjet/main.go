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
package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/crossplaneio/crossplane-tools/internal/comments"
	"github.com/crossplaneio/crossplane-tools/internal/fields"
	"github.com/crossplaneio/crossplane-tools/internal/generate"
	"github.com/crossplaneio/crossplane-tools/internal/match"
	"github.com/crossplaneio/crossplane-tools/internal/method"
)

const (
	// LoadMode used to load all packages.
	LoadMode = packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedDeps | packages.NeedTypes | packages.NeedSyntax

	// DisableMarker used to disable generation of managed resource methods for
	// a type that otherwise appears to be a managed resource that is missing a
	// subnet of its methods.
	DisableMarker = "crossplane:generate:methods"
)

// Imports used in generated code.
const (
	CoreAlias  = "corev1"
	CoreImport = "k8s.io/api/core/v1"

	MetaAlias  = "metav1"
	MetaImport = "k8s.io/apimachinery/pkg/apis/meta/v1"

	RuntimeAlias  = "runtimev1alpha1"
	RuntimeImport = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
)

func main() {
	var (
		app = kingpin.New(filepath.Base(os.Args[0]), "Generates Crossplane API type methods.").DefaultEnvars()

		methodsets      = app.Command("generate-methodsets", "Generate a Crossplane method sets.")
		headerFile      = methodsets.Flag("header-file", "The contents of this file will be added to the top of all generated files.").ExistingFile()
		filenameManaged = methodsets.Flag("filename-managed", "The filename of generated managed resource files.").Default("zz_generated.managed.go").String()
		filenameClaim   = methodsets.Flag("filename-claim", "The filename of generated resource claim files.").Default("zz_generated.claim.go").String()
		filenameClass   = methodsets.Flag("filename-class", "The filename of generated resource class files.").Default("zz_generated.class.go").String()
		pattern         = methodsets.Arg("packages", "Package(s) for which to generate methods, for example github.com/crossplaneio/crossplane/apis/...").String()
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	pkgs, err := packages.Load(&packages.Config{Mode: LoadMode}, *pattern)
	kingpin.FatalIfError(err, "cannot load packages %s", *pattern)

	header := ""
	if *headerFile != "" {
		h, err := ioutil.ReadFile(*headerFile)
		kingpin.FatalIfError(err, "cannot read header file %s", *headerFile)
		header = string(h)
	}

	for _, p := range pkgs {
		for _, err := range p.Errors {
			kingpin.FatalIfError(err, "error loading packages using pattern %s", *pattern)
		}
		kingpin.FatalIfError(GenerateManaged(*filenameManaged, header, p), "cannot write managed resource method set for package %s", p.PkgPath)
		kingpin.FatalIfError(GenerateClaim(*filenameClaim, header, p), "cannot write resource claim method set for package %s", p.PkgPath)
		kingpin.FatalIfError(GenerateClass(*filenameClass, header, p), "cannot write resource class method set for package %s", p.PkgPath)
	}
}

// GenerateManaged generates the resource.Managed method set.
func GenerateManaged(filename, header string, p *packages.Package) error {
	receiver := "mg"

	methods := method.Set{
		"SetConditions":                       method.NewSetConditions(receiver, RuntimeImport),
		"SetBindingPhase":                     method.NewSetBindingPhase(receiver, RuntimeImport),
		"GetBindingPhase":                     method.NewGetBindingPhase(receiver, RuntimeImport),
		"SetClaimReference":                   method.NewSetClaimReference(receiver, CoreImport),
		"GetClaimReference":                   method.NewGetClaimReference(receiver, CoreImport),
		"SetClassReference":                   method.NewSetClassReference(receiver, CoreImport),
		"GetClassReference":                   method.NewGetClassReference(receiver, CoreImport),
		"SetWriteConnectionSecretToReference": method.NewSetWriteConnectionSecretToReference(receiver, RuntimeImport),
		"GetWriteConnectionSecretToReference": method.NewGetWriteConnectionSecretToReference(receiver, RuntimeImport),
		"SetReclaimPolicy":                    method.NewSetReclaimPolicy(receiver, RuntimeImport, fields.NameSpec),
		"GetReclaimPolicy":                    method.NewGetReclaimPolicy(receiver, RuntimeImport, fields.NameSpec),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
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

// GenerateClaim generates the resource.Claim method set.
func GenerateClaim(filename, header string, p *packages.Package) error {
	receiver := "cm"

	methods := method.Set{
		"SetConditions":                       method.NewSetConditions(receiver, RuntimeImport),
		"SetBindingPhase":                     method.NewSetBindingPhase(receiver, RuntimeImport),
		"GetBindingPhase":                     method.NewGetBindingPhase(receiver, RuntimeImport),
		"SetClassSelector":                    method.NewSetClassSelector(receiver, MetaImport),
		"GetClassSelector":                    method.NewGetClassSelector(receiver, MetaImport),
		"SetClassReference":                   method.NewSetClassReference(receiver, CoreImport),
		"GetClassReference":                   method.NewGetClassReference(receiver, CoreImport),
		"SetResourceReference":                method.NewSetResourceReference(receiver, CoreImport),
		"GetResourceReference":                method.NewGetResourceReference(receiver, CoreImport),
		"SetWriteConnectionSecretToReference": method.NewLocalSetWriteConnectionSecretToReference(receiver, RuntimeImport),
		"GetWriteConnectionSecretToReference": method.NewLocalGetWriteConnectionSecretToReference(receiver, RuntimeImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{
			CoreImport:    CoreAlias,
			MetaImport:    MetaAlias,
			RuntimeImport: RuntimeAlias,
		}),
		generate.WithMatcher(match.AllOf(
			match.Claim(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, "cannot write resource claim methods")
}

// GenerateClass generates the resource.Class method set.
func GenerateClass(filename, header string, p *packages.Package) error {
	receiver := "cs"

	methods := method.Set{
		"SetReclaimPolicy": method.NewSetReclaimPolicy(receiver, RuntimeImport, fields.NameSpecTemplate),
		"GetReclaimPolicy": method.NewGetReclaimPolicy(receiver, RuntimeImport, fields.NameSpecTemplate),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{RuntimeImport: RuntimeAlias}),
		generate.WithMatcher(match.AllOf(
			match.Class(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, "cannot write resource class methods")
}
