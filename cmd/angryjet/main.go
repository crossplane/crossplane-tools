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
	"github.com/crossplane/crossplane-tools/internal/generate"
	"github.com/crossplane/crossplane-tools/internal/match"
	"github.com/crossplane/crossplane-tools/internal/method"
	"github.com/crossplane/crossplane-tools/internal/types"
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

	ClientAlias  = "client"
	ClientImport = "sigs.k8s.io/controller-runtime/pkg/client"

	RuntimeAlias  = "xpv1"
	RuntimeImport = "github.com/crossplane/crossplane-runtime/apis/common/v1"

	ResourceAlias  = "resource"
	ResourceImport = "github.com/crossplane/crossplane-runtime/pkg/resource"

	ReferenceAlias  = "reference"
	ReferenceImport = "github.com/crossplane/crossplane-runtime/pkg/reference"
)

func main() {
	var (
		app = kingpin.New(filepath.Base(os.Args[0]), "Generates Crossplane API type methods.").DefaultEnvars()

		methodsets          = app.Command("generate-methodsets", "Generate a Crossplane method sets.")
		headerFile          = methodsets.Flag("header-file", "The contents of this file will be added to the top of all generated files.").ExistingFile()
		filenameManaged     = methodsets.Flag("filename-managed", "The filename of generated managed resource files.").Default("zz_generated.managed.go").String()
		filenameResolvers   = methodsets.Flag("filename-resolvers", "The filename of generated reference resolver files.").Default("zz_generated.resolvers.go").String()
		filenameManagedList = methodsets.Flag("filename-managed-list", "The filename of generated managed list resource files.").Default("zz_generated.managedlist.go").String()
		filenamePC          = methodsets.Flag("filename-pc", "The filename of generated provider config files.").Default("zz_generated.pc.go").String()
		filenamePCU         = methodsets.Flag("filename-pcu", "The filename of generated provider config usage files.").Default("zz_generated.pcu.go").String()
		filenamePCUList     = methodsets.Flag("filename-pcu-list", "The filename of generated provider config usage files.").Default("zz_generated.pculist.go").String()
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
		kingpin.FatalIfError(GenerateProviderConfig(*filenamePC, header, p), "cannot write provider config method set for package %s", p.PkgPath)
		kingpin.FatalIfError(GenerateProviderConfigUsage(*filenamePCU, header, p), "cannot write provider config usage method set for package %s", p.PkgPath)
		kingpin.FatalIfError(GenerateProviderConfigUsageList(*filenamePCUList, header, p), "cannot write provider config usage list method set for package %s", p.PkgPath)
		kingpin.FatalIfError(GenerateReferences(*filenameResolvers, header, p), "cannot write reference resolvers for package %s", p.PkgPath)
	}
}

// GenerateManaged generates the resource.Managed method set.
func GenerateManaged(filename, header string, p *packages.Package) error {
	receiver := "mg"

	methods := method.Set{
		"SetConditions":                       method.NewSetConditions(receiver, RuntimeImport),
		"GetCondition":                        method.NewGetCondition(receiver, RuntimeImport),
		"GetProviderReference":                method.NewGetProviderReference(receiver, RuntimeImport),
		"SetProviderReference":                method.NewSetProviderReference(receiver, RuntimeImport),
		"GetProviderConfigReference":          method.NewGetProviderConfigReference(receiver, RuntimeImport),
		"SetProviderConfigReference":          method.NewSetProviderConfigReference(receiver, RuntimeImport),
		"SetWriteConnectionSecretToReference": method.NewSetWriteConnectionSecretToReference(receiver, RuntimeImport),
		"GetWriteConnectionSecretToReference": method.NewGetWriteConnectionSecretToReference(receiver, RuntimeImport),
		"SetPublishConnectionDetailsTo":       method.NewSetPublishConnectionDetailsTo(receiver, RuntimeImport),
		"GetPublishConnectionDetailsTo":       method.NewGetPublishConnectionDetailsTo(receiver, RuntimeImport),
		"SetManagementPolicy":                 method.NewSetManagementPolicy(receiver, RuntimeImport),
		"GetManagementPolicy":                 method.NewGetManagementPolicy(receiver, RuntimeImport),
		"SetDeletionPolicy":                   method.NewSetDeletionPolicy(receiver, RuntimeImport),
		"GetDeletionPolicy":                   method.NewGetDeletionPolicy(receiver, RuntimeImport),
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

// GenerateProviderConfig generates the resource.ProviderConfig method set.
func GenerateProviderConfig(filename, header string, p *packages.Package) error {
	receiver := "p"

	methods := method.Set{
		"SetUsers":      method.NewSetUsers(receiver),
		"GetUsers":      method.NewGetUsers(receiver),
		"SetConditions": method.NewSetConditions(receiver, RuntimeImport),
		"GetCondition":  method.NewGetCondition(receiver, RuntimeImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{RuntimeImport: RuntimeAlias}),
		generate.WithMatcher(match.AllOf(
			match.ProviderConfig(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, "cannot write provider config methods")
}

// GenerateProviderConfigUsage generates the resource.ProviderConfigUsage method set.
func GenerateProviderConfigUsage(filename, header string, p *packages.Package) error {
	receiver := "p"

	methods := method.Set{
		"SetProviderConfigReference": method.NewSetRootProviderConfigReference(receiver, RuntimeImport),
		"GetProviderConfigReference": method.NewGetRootProviderConfigReference(receiver, RuntimeImport),
		"SetResourceReference":       method.NewSetRootResourceReference(receiver, RuntimeImport),
		"GetResourceReference":       method.NewGetRootResourceReference(receiver, RuntimeImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{RuntimeImport: RuntimeAlias}),
		generate.WithMatcher(match.AllOf(
			match.ProviderConfigUsage(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, "cannot write provider config usage methods")
}

// GenerateProviderConfigUsageList generates the
// resource.ProviderConfigUsageList method set.
func GenerateProviderConfigUsageList(filename, header string, p *packages.Package) error {
	receiver := "p"

	methods := method.Set{
		"GetItems": method.NewProviderConfigUsageGetItems(receiver, ResourceImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{RuntimeImport: RuntimeAlias}),
		generate.WithMatcher(match.AllOf(
			match.ProviderConfigUsageList(),
			match.DoesNotHaveMarker(comments.In(p), DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, "cannot write provider config usage list methods")
}

// GenerateReferences generates reference resolver calls.
func GenerateReferences(filename, header string, p *packages.Package) error {
	receiver := "mg"
	comm := comments.In(p)

	methods := method.Set{
		"ResolveReferences": method.NewResolveReferences(types.NewTraverser(comm), receiver, ClientImport, ReferenceImport),
	}

	err := generate.WriteMethods(p, methods, filepath.Join(filepath.Dir(p.GoFiles[0]), filename),
		generate.WithHeaders(header),
		generate.WithImportAliases(map[string]string{
			ClientImport:    ClientAlias,
			ReferenceImport: ReferenceAlias,
		}),
		generate.WithMatcher(match.AllOf(
			match.Managed(),
			match.DoesNotHaveMarker(comm, DisableMarker, "false")),
		),
	)

	return errors.Wrap(err, "cannot write reference resolver methods")
}
