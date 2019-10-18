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

	RuntimeAlias  = "runtimev1alpha1"
	RuntimeImport = "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"

	ResourceImport = "github.com/crossplaneio/crossplane-runtime/pkg/resource"
)

func main() {
	var (
		app = kingpin.New(filepath.Base(os.Args[0]), "Generates Crossplane API type methods.").DefaultEnvars()

		methodsets                = app.Command("generate-methodsets", "Generate a Crossplane method sets.")
		headerFile                = methodsets.Flag("header-file", "The contents of this file will be added to the top of all generated files.").ExistingFile()
		filenameManaged           = methodsets.Flag("filename-managed", "The filename of generated managed resource files.").Default("zz_generated.managed.go").String()
		filenameClaim             = methodsets.Flag("filename-claim", "The filename of generated resource claim files.").Default("zz_generated.claim.go").String()
		filenamePortableClass     = methodsets.Flag("filename-portable-class", "The filename of generated portable class files.").Default("zz_generated.portableclass.go").String()
		filenamePortableClassList = methodsets.Flag("filename-portable-class-list", "The filename of generated portable class list files.").Default("zz_generated.portableclasslist.go").String()
		filenameNonPortableClass  = methodsets.Flag("filename-non-portable-class", "The filename of generated non-portable class files.").Default("zz_generated.nonportableclass.go").String()
		pattern                   = methodsets.Arg("packages", "Package(s) for which to generate methods, for example github.com/crossplaneio/crossplane/apis/...").String()
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
		kingpin.FatalIfError(GeneratePortableClass(*filenamePortableClass, header, p), "cannot write portable class method set for package %s", p.PkgPath)
		kingpin.FatalIfError(GeneratePortableClassList(*filenamePortableClassList, header, p), "cannot write portable class list method set for package %s", p.PkgPath)
		kingpin.FatalIfError(GenerateNonPortableClass(*filenameNonPortableClass, header, p), "cannot write non portable class method set for package %s", p.PkgPath)
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
		"SetNonPortableClassReference":        method.NewSetNonPortableClassReference(receiver, CoreImport),
		"GetNonPortableClassReference":        method.NewGetNonPortableClassReference(receiver, CoreImport),
		"SetWriteConnectionSecretToReference": method.NewSetWriteConnectionSecretToReference(receiver, CoreImport),
		"GetWriteConnectionSecretToReference": method.NewGetWriteConnectionSecretToReference(receiver, CoreImport),
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
		"SetResourceReference":                method.NewSetResourceReference(receiver, CoreImport),
		"GetResourceReference":                method.NewGetResourceReference(receiver, CoreImport),
		"SetPortableClassReference":           method.NewSetPortableClassReference(receiver, CoreImport),
		"GetPortableClassReference":           method.NewGetPortableClassReference(receiver, CoreImport),
		"SetWriteConnectionSecretToReference": method.NewSetWriteConnectionSecretToReference(receiver, CoreImport),
		"GetWriteConnectionSecretToReference": method.NewGetWriteConnectionSecretToReference(receiver, CoreImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{
			CoreImport:    CoreAlias,
			RuntimeImport: RuntimeAlias,
		}),
		generate.WithMatcher(match.AllOf(
			match.Claim(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, "cannot write resource claim methods")
}

// GeneratePortableClass generates the resource.PortableClass method set.
func GeneratePortableClass(filename, header string, p *packages.Package) error {
	receiver := "cs"

	methods := method.Set{
		"SetNonPortableClassReference": method.NewSetNonPortableClassReference(receiver, CoreImport),
		"GetNonPortableClassReference": method.NewGetNonPortableClassReference(receiver, CoreImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{CoreImport: CoreAlias}),
		generate.WithMatcher(match.AllOf(
			match.PortableClass(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, "cannot write portable class methods")
}

// GeneratePortableClassList generates the resource.PortableClassList method set.
func GeneratePortableClassList(filename, header string, p *packages.Package) error {
	receiver := "csl"

	methods := method.Set{
		"SetPortableClassItems": method.NewSetPortableClassItems(receiver, ResourceImport),
		"GetPortableClassItems": method.NewGetPortableClassItems(receiver, ResourceImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithMatcher(match.AllOf(
			match.PortableClassList(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, "cannot write portable class methods")
}

// GenerateNonPortableClass generates the resource.NonPortableClass method set.
func GenerateNonPortableClass(filename, header string, p *packages.Package) error {
	receiver := "cs"

	methods := method.Set{
		"SetReclaimPolicy": method.NewSetReclaimPolicy(receiver, RuntimeImport, fields.NameSpecTemplate),
		"GetReclaimPolicy": method.NewGetReclaimPolicy(receiver, RuntimeImport, fields.NameSpecTemplate),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{RuntimeImport: RuntimeAlias}),
		generate.WithMatcher(match.AllOf(
			match.NonPortableClass(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, "cannot write non-portable class methods")
}
