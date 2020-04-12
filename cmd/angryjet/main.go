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

	"github.com/crossplane/crossplane-tools/internal/comments"
	"github.com/crossplane/crossplane-tools/internal/fields"
	"github.com/crossplane/crossplane-tools/internal/generate"
	"github.com/crossplane/crossplane-tools/internal/match"
	"github.com/crossplane/crossplane-tools/internal/method"
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
	RuntimeImport = "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"

	ResourceAlias  = "resource"
	ResourceImport = "github.com/crossplane/crossplane-runtime/pkg/resource"
)

func main() {
	var (
		app = kingpin.New(filepath.Base(os.Args[0]), "Generates Crossplane API type methods.").DefaultEnvars()

		methodsets          = app.Command("generate-methodsets", "Generate a Crossplane method sets.")
		headerFile          = methodsets.Flag("header-file", "The contents of this file will be added to the top of all generated files.").ExistingFile()
		filenameManaged     = methodsets.Flag("filename-managed", "The filename of generated managed resource files.").Default("zz_generated.managed.go").String()
		filenameManagedList = methodsets.Flag("filename-managed-list", "The filename of generated managed list resource files.").Default("zz_generated.managedlist.go").String()
		filenameClaim       = methodsets.Flag("filename-claim", "The filename of generated resource claim files.").Default("zz_generated.claim.go").String()
		filenameClaimList   = methodsets.Flag("filename-claim-list", "The filename of generated resource claim list files.").Default("zz_generated.claimlist.go").String()
		filenameClass       = methodsets.Flag("filename-class", "The filename of generated resource class files.").Default("zz_generated.class.go").String()
		filenameClassList   = methodsets.Flag("filename-class-list", "The filename of generated resource class list files.").Default("zz_generated.classlist.go").String()
		filenameProvider    = methodsets.Flag("filename-provider", "The filename of generated provider files.").Default("zz_generated.provider.go").String()
		pattern             = methodsets.Arg("packages", "Package(s) for which to generate methods, for example github.com/crossplane/crossplane/apis/...").String()
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
		kingpin.FatalIfError(GenerateManagedList(*filenameManagedList, header, p), "cannot write managed resource list method set for package %s", p.PkgPath)
		kingpin.FatalIfError(GenerateClaim(*filenameClaim, header, p), "cannot write resource claim method set for package %s", p.PkgPath)
		kingpin.FatalIfError(GenerateClaimList(*filenameClaimList, header, p), "cannot write resource claim list method set for package %s", p.PkgPath)
		kingpin.FatalIfError(GenerateClass(*filenameClass, header, p), "cannot write resource class method set for package %s", p.PkgPath)
		kingpin.FatalIfError(GenerateClassList(*filenameClassList, header, p), "cannot write resource class list method set for package %s", p.PkgPath)
		kingpin.FatalIfError(GenerateProvider(*filenameProvider, header, p), "cannot write provider method set for package %s", p.PkgPath)
	}
}

// GenerateManaged generates the resource.Managed method set.
func GenerateManaged(filename, header string, p *packages.Package) error {
	receiver := "mg"

	methods := method.Set{
		"SetConditions":                       method.NewSetConditions(receiver, RuntimeImport),
		"GetCondition":                        method.NewGetCondition(receiver, RuntimeImport),
		"SetBindingPhase":                     method.NewSetBindingPhase(receiver, RuntimeImport),
		"GetBindingPhase":                     method.NewGetBindingPhase(receiver, RuntimeImport),
		"SetClaimReference":                   method.NewSetClaimReference(receiver, CoreImport),
		"GetClaimReference":                   method.NewGetClaimReference(receiver, CoreImport),
		"SetClassReference":                   method.NewSetClassReference(receiver, CoreImport),
		"GetClassReference":                   method.NewGetClassReference(receiver, CoreImport),
		"GetProviderReference":                method.NewGetProviderReference(receiver, CoreImport),
		"SetProviderReference":                method.NewSetProviderReference(receiver, CoreImport),
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

// GenerateManagedList generates the resource.ManagedList method set.
func GenerateManagedList(filename, header string, p *packages.Package) error {
	receiver := "l"

	methods := method.Set{
		"GetItems": method.NewManagedGetItems(receiver, ResourceImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{
			ResourceImport: ResourceAlias,
		}),
		generate.WithMatcher(match.AllOf(
			match.ManagedList(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, "cannot write managed resource list methods")
}

// GenerateClaim generates the resource.Claim method set.
func GenerateClaim(filename, header string, p *packages.Package) error {
	receiver := "cm"

	methods := method.Set{
		"SetConditions":                       method.NewSetConditions(receiver, RuntimeImport),
		"GetCondition":                        method.NewGetCondition(receiver, RuntimeImport),
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

// GenerateClaimList generates the resource.ClaimList method set.
func GenerateClaimList(filename, header string, p *packages.Package) error {
	receiver := "l"

	methods := method.Set{
		"GetItems": method.NewClaimGetItems(receiver, ResourceImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{
			ResourceImport: ResourceAlias,
		}),
		generate.WithMatcher(match.AllOf(
			match.ClaimList(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, "cannot write resource claim list methods")
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

// GenerateClassList generates the resource.ClassList method set.
func GenerateClassList(filename, header string, p *packages.Package) error {
	receiver := "l"

	methods := method.Set{
		"GetItems": method.NewClassGetItems(receiver, ResourceImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{
			ResourceImport: ResourceAlias,
		}),
		generate.WithMatcher(match.AllOf(
			match.ClassList(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, "cannot write resource class list methods")
}

// GenerateProvider generates the resource.Provider method set.
func GenerateProvider(filename, header string, p *packages.Package) error {
	receiver := "p"

	methods := method.Set{
		"SetCredentialsSecretReference": method.NewSetCredentialsSecretReference(receiver, RuntimeImport),
		"GetCredentialsSecretReference": method.NewGetCredentialsSecretReference(receiver, RuntimeImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{RuntimeImport: RuntimeAlias}),
		generate.WithMatcher(match.AllOf(
			match.Provider(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, "cannot write provider methods")
}
