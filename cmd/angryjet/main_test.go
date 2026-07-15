/*
Copyright 2026 The Crossplane Authors.

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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/packages/packagestest"
)

// The fixture below models one kind for each shape angryjet generates for:
//
//   - Legacy/Modern: the crossplane-runtime common/v1 (cluster) and common/v2
//     (namespaced) types the generators have always supported.
//   - Core: the same shapes but sourced from the crossplane/apis/v2/core/v2
//     module that Crossplane v2.3 moved the common types into. These are matched
//     by the *Core generators.
const (
	angryjetFixtureSource = `package v1alpha1

import (
	xprv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	xprv2 "github.com/crossplane/crossplane-runtime/v2/apis/common/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

type ReferenceTargetSpec struct {
	xpv2.ClusterManagedResourceSpec
}

type ReferenceTargetStatus struct {
	xpv2.ManagedResourceStatus
}

type ReferenceTarget struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   ReferenceTargetSpec
	Status ReferenceTargetStatus
}

type ReferenceTargetList struct {
	metav1.TypeMeta
	Items []ReferenceTarget
}

type NamespacedResourceParameters struct {
	// +crossplane:generate:reference:type=ReferenceTarget
	Target string
}

type NamespacedResourceSpec struct {
	xpv2.ManagedResourceSpec
	ForProvider NamespacedResourceParameters
}

type NamespacedResourceStatus struct {
	xpv2.ManagedResourceStatus
}

type NamespacedResource struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   NamespacedResourceSpec
	Status NamespacedResourceStatus
}

type NamespacedResourceList struct {
	metav1.TypeMeta
	Items []NamespacedResource
}

type ClusterResourceParameters struct {
	// +crossplane:generate:reference:type=ReferenceTarget
	Target string
}

type ClusterResourceSpec struct {
	xpv2.ClusterManagedResourceSpec
	ForProvider ClusterResourceParameters
}

type ClusterResourceStatus struct {
	xpv2.ManagedResourceStatus
}

type ClusterResource struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   ClusterResourceSpec
	Status ClusterResourceStatus
}

type ClusterResourceList struct {
	metav1.TypeMeta
	Items []ClusterResource
}

type ModernResourceParameters struct {
	// +crossplane:generate:reference:type=ReferenceTarget
	Target string
}

type ModernResourceSpec struct {
	xprv2.ManagedResourceSpec
	ForProvider ModernResourceParameters
}

type ModernResourceStatus struct {
	xprv1.ResourceStatus
}

type ModernResource struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   ModernResourceSpec
	Status ModernResourceStatus
}

type ModernResourceList struct {
	metav1.TypeMeta
	Items []ModernResource
}

// A legacy (cluster-scoped) managed resource using crossplane-runtime common/v1
// types, exercising GenerateManagedLegacy, GenerateManagedListLegacy, and the
// non-namespaced GenerateReferencesLegacy.
type LegacyResourceParameters struct {
	// +crossplane:generate:reference:type=ReferenceTarget
	Target string
}

type LegacyResourceSpec struct {
	xprv1.ResourceSpec
	ForProvider LegacyResourceParameters
}

type LegacyResourceStatus struct {
	xprv1.ResourceStatus
}

type LegacyResource struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   LegacyResourceSpec
	Status LegacyResourceStatus
}

type LegacyResourceList struct {
	metav1.TypeMeta
	Items []LegacyResource
}

// A provider config usage embedding the crossplane-runtime common/v1 non-typed
// usage, exercising GenerateProviderConfigUsageLegacy and
// GenerateProviderConfigUsageListLegacy.
type RuntimeLegacyProviderConfigUsage struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	xprv1.ProviderConfigUsage
}

type RuntimeLegacyProviderConfigUsageList struct {
	metav1.TypeMeta
	Items []RuntimeLegacyProviderConfigUsage
}

// A provider config usage embedding the crossplane-runtime common/v2 typed
// usage, exercising the common/v2 GenerateProviderConfigUsageModern.
type RuntimeModernProviderConfigUsage struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	xprv2.TypedProviderConfigUsage
}

type RuntimeModernProviderConfigUsageList struct {
	metav1.TypeMeta
	Items []RuntimeModernProviderConfigUsage
}

// A provider config usage embedding the core API v2 typed usage, exercising
// GenerateProviderConfigUsageModernCore.
type CoreModernProviderConfigUsage struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	xpv2.TypedProviderConfigUsage
}

type CoreModernProviderConfigUsageList struct {
	metav1.TypeMeta
	Items []CoreModernProviderConfigUsage
}

// A provider config usage embedding the core API v2 non-typed usage, exercising
// GenerateProviderConfigUsageLegacyCore.
type CoreLegacyProviderConfigUsage struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	xpv2.ProviderConfigUsage
}

type CoreLegacyProviderConfigUsageList struct {
	metav1.TypeMeta
	Items []CoreLegacyProviderConfigUsage
}

type CoreProviderCredentials struct {
	Source xpv2.CredentialsSource
	xpv2.CommonCredentialSelectors
}

type CoreProviderConfigSpec struct {
	xprv1.ProviderConfigSpec
	Credentials CoreProviderCredentials
}

// A CoreProviderConfigStatus reflects the observed state of a ProviderConfig.
type CoreProviderConfigStatus struct {
	xpv2.ProviderConfigStatus
}

type CoreProviderConfig struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   CoreProviderConfigSpec
	Status CoreProviderConfigStatus
}

type LegacyProviderConfigSpec struct {
	xprv1.ProviderConfigSpec
}

type LegacyProviderConfigStatus struct {
	xprv1.ProviderConfigStatus
}

type LegacyProviderConfig struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   LegacyProviderConfigSpec
	Status LegacyProviderConfigStatus
}
`

	metaV1FixtureSource = `package v1

type TypeMeta struct{}

type ObjectMeta struct{}

type ListMeta struct{}
`

	coreV2FixtureSource = `package v2

type ManagedResourceSpec struct{}

type ClusterManagedResourceSpec struct{}

type ManagedResourceStatus struct{}
`

	runtimeV1FixtureSource = `package v1

type ResourceSpec struct{}

type ResourceStatus struct{}

type Condition struct{}

type ConditionType string

type User struct{}

type ProviderConfigSpec struct{}

type ProviderConfigUsage struct{}

type ProviderConfigStatus struct{}

func (p *ProviderConfigStatus) SetConditions(c ...Condition) {}

func (p *ProviderConfigStatus) GetCondition(ct ConditionType) Condition { return Condition{} }

func (p *ProviderConfigStatus) SetUsers(u ...User) {}

func (p *ProviderConfigStatus) GetUsers() []User { return nil }
`

	runtimeV2FixtureSource = `package v2

type ManagedResourceSpec struct{}

type TypedProviderConfigUsage struct{}
`

	coreV2ProviderConfigFixtureSource = `package v2

type ProviderConfigStatus struct{}

type ProviderConfigReference struct{}

type Reference struct{}

type TypedReference struct{}

type ProviderConfigUsage struct{}

type TypedProviderConfigUsage struct{}

type CredentialsSource string

type CommonCredentialSelectors struct{}
`
)

type generatorFunc func(string, string, *packages.Package) error

func TestGenerateMethods(t *testing.T) {
	type args struct {
		generate generatorFunc
		filename string
	}

	type want struct {
		contains    []string
		notContains []string
	}

	cases := map[string]struct {
		args
		want
	}{
		"Managed modern core": {
			args: args{
				generate: GenerateManagedModernCore,
				filename: "zz_generated.modern_core.go",
			},
			want: want{
				contains: []string{
					`import xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"`,
					`func (mg *NamespacedResource) SetProviderConfigReference(r *xpv2.ProviderConfigReference) {`,
					`func (mg *NamespacedResource) GetWriteConnectionSecretToReference() *xpv2.LocalSecretReference {`,
					`func (mg *NamespacedResource) GetManagementPolicies() xpv2.ManagementPolicies {`,
				},
				notContains: []string{
					`func (mg *ClusterResource)`,
					`SetDeletionPolicy`,
				},
			},
		},
		"Managed legacy core": {
			args: args{
				generate: GenerateManagedLegacyCore,
				filename: "zz_generated.legacy_core.go",
			},
			want: want{
				contains: []string{
					`import xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"`,
					`func (mg *ClusterResource) SetProviderConfigReference(r *xpv2.Reference) {`,
					`func (mg *ClusterResource) GetWriteConnectionSecretToReference() *xpv2.SecretReference {`,
					`func (mg *ClusterResource) SetDeletionPolicy(r xpv2.DeletionPolicy) {`,
				},
				notContains: []string{
					`func (mg *NamespacedResource)`,
				},
			},
		},
		"Managed list modern core": {
			args: args{
				generate: GenerateManagedListModernCore,
				filename: "zz_generated.modern_core_list.go",
			},
			want: want{
				contains: []string{
					`import resource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"`,
					`func (l *NamespacedResourceList) GetItems() []resource.Managed {`,
				},
				notContains: []string{
					`func (l *ClusterResourceList)`,
				},
			},
		},
		"Managed list legacy core": {
			args: args{
				generate: GenerateManagedListLegacyCore,
				filename: "zz_generated.legacy_core_list.go",
			},
			want: want{
				contains: []string{
					`import resource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"`,
					`func (l *ClusterResourceList) GetItems() []resource.Managed {`,
				},
				notContains: []string{
					`func (l *NamespacedResourceList)`,
				},
			},
		},
		"References modern core": {
			args: args{
				generate: GenerateReferencesModernCore,
				filename: "zz_generated.modern_core_refs.go",
			},
			want: want{
				contains: []string{
					`reference.NewAPINamespacedResolver(c, mg)`,
					`var rsp reference.NamespacedResolutionResponse`,
					`reference.NamespacedResolutionRequest{`,
					`Namespace:    mg.GetNamespace(),`,
					`func (mg *NamespacedResource) ResolveReferences(ctx context.Context, c client.Reader) error {`,
				},
			},
		},
		"References legacy core": {
			args: args{
				generate: GenerateReferencesLegacyCore,
				filename: "zz_generated.legacy_core_refs.go",
			},
			want: want{
				contains: []string{
					`reference.NewAPIResolver(c, mg)`,
					`var rsp reference.ResolutionResponse`,
					`reference.ResolutionRequest{`,
					`Namespace:    mg.GetNamespace(),`,
					`func (mg *ClusterResource) ResolveReferences(ctx context.Context, c client.Reader) error {`,
				},
			},
		},
		"Managed modern": {
			args: args{
				generate: GenerateManagedModern,
				filename: "zz_generated.modern.go",
			},
			want: want{
				contains: []string{
					`import xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"`,
					`func (mg *ModernResource) SetProviderConfigReference(r *xpv1.ProviderConfigReference) {`,
					`func (mg *ModernResource) GetWriteConnectionSecretToReference() *xpv1.LocalSecretReference {`,
				},
				notContains: []string{
					`func (mg *NamespacedResource)`,
					`SetDeletionPolicy`,
				},
			},
		},
		"Managed legacy": {
			args: args{
				generate: GenerateManagedLegacy,
				filename: "zz_generated.legacy.go",
			},
			want: want{
				contains: []string{
					`import xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"`,
					`func (mg *LegacyResource) SetProviderConfigReference(r *xpv1.Reference) {`,
					`func (mg *LegacyResource) GetWriteConnectionSecretToReference() *xpv1.SecretReference {`,
					`func (mg *LegacyResource) SetDeletionPolicy(r xpv1.DeletionPolicy) {`,
				},
				notContains: []string{
					`func (mg *ModernResource)`,
					`func (mg *ClusterResource)`,
				},
			},
		},
		"Managed list modern": {
			args: args{
				generate: GenerateManagedListModern,
				filename: "zz_generated.modern_list.go",
			},
			want: want{
				contains: []string{
					`func (l *ModernResourceList) GetItems() []resource.Managed {`,
				},
				notContains: []string{
					`func (l *NamespacedResourceList)`,
				},
			},
		},
		"Managed list legacy": {
			args: args{
				generate: GenerateManagedListLegacy,
				filename: "zz_generated.legacy_list.go",
			},
			want: want{
				contains: []string{
					`import resource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"`,
					`func (l *LegacyResourceList) GetItems() []resource.Managed {`,
				},
				notContains: []string{
					`func (l *ModernResourceList)`,
				},
			},
		},
		"Provider config usage modern": {
			args: args{
				generate: GenerateProviderConfigUsageModern,
				filename: "zz_generated.modern_pcu.go",
			},
			want: want{
				contains: []string{
					`func (p *RuntimeModernProviderConfigUsage) SetProviderConfigReference(r xpv1.ProviderConfigReference) {`,
					`func (p *RuntimeModernProviderConfigUsage) GetProviderConfigReference() xpv1.ProviderConfigReference {`,
					`func (p *RuntimeModernProviderConfigUsage) SetResourceReference(r xpv1.TypedReference) {`,
					`func (p *RuntimeModernProviderConfigUsage) GetResourceReference() xpv1.TypedReference {`,
				},
				notContains: []string{
					`func (p *RuntimeModernProviderConfigUsage) SetProviderConfigReference(r xpv2.ProviderConfigReference) {`,
					`func (p *CoreModernProviderConfigUsage)`,
				},
			},
		},
		"Provider config usage legacy": {
			args: args{
				generate: GenerateProviderConfigUsageLegacy,
				filename: "zz_generated.legacy_pcu.go",
			},
			want: want{
				contains: []string{
					`func (p *RuntimeLegacyProviderConfigUsage) SetProviderConfigReference(r xpv1.Reference) {`,
					`func (p *RuntimeLegacyProviderConfigUsage) GetProviderConfigReference() xpv1.Reference {`,
					`func (p *RuntimeLegacyProviderConfigUsage) SetResourceReference(r xpv1.TypedReference) {`,
					`func (p *RuntimeLegacyProviderConfigUsage) GetResourceReference() xpv1.TypedReference {`,
				},
				notContains: []string{
					`func (p *RuntimeLegacyProviderConfigUsage) SetProviderConfigReference(r xpv1.ProviderConfigReference) {`,
					`func (p *CoreLegacyProviderConfigUsage)`,
				},
			},
		},
		"Provider config usage modern core": {
			args: args{
				generate: GenerateProviderConfigUsageModernCore,
				filename: "zz_generated.modern_core_pcu.go",
			},
			want: want{
				contains: []string{
					`func (p *CoreModernProviderConfigUsage) SetProviderConfigReference(r xpv2.ProviderConfigReference) {`,
					`func (p *CoreModernProviderConfigUsage) GetProviderConfigReference() xpv2.ProviderConfigReference {`,
					`func (p *CoreModernProviderConfigUsage) SetResourceReference(r xpv2.TypedReference) {`,
					`func (p *CoreModernProviderConfigUsage) GetResourceReference() xpv2.TypedReference {`,
				},
				notContains: []string{
					`func (p *CoreModernProviderConfigUsage) SetProviderConfigReference(r xpv1.ProviderConfigReference) {`,
					`func (p *RuntimeModernProviderConfigUsage)`,
				},
			},
		},
		"Provider config usage legacy core": {
			args: args{
				generate: GenerateProviderConfigUsageLegacyCore,
				filename: "zz_generated.legacy_core_pcu.go",
			},
			want: want{
				contains: []string{
					`func (p *CoreLegacyProviderConfigUsage) SetProviderConfigReference(r xpv2.Reference) {`,
					`func (p *CoreLegacyProviderConfigUsage) GetProviderConfigReference() xpv2.Reference {`,
					`func (p *CoreLegacyProviderConfigUsage) SetResourceReference(r xpv2.TypedReference) {`,
					`func (p *CoreLegacyProviderConfigUsage) GetResourceReference() xpv2.TypedReference {`,
				},
				notContains: []string{
					`func (p *CoreLegacyProviderConfigUsage) SetProviderConfigReference(r xpv2.ProviderConfigReference) {`,
					`func (p *CoreModernProviderConfigUsage)`,
				},
			},
		},
		"Provider config usage list modern": {
			args: args{
				generate: GenerateProviderConfigUsageListModern,
				filename: "zz_generated.modern_pcu_list.go",
			},
			want: want{
				contains: []string{
					`func (p *RuntimeModernProviderConfigUsageList) GetItems() []resource.ProviderConfigUsage {`,
				},
				notContains: []string{
					`func (p *CoreModernProviderConfigUsageList)`,
					`func (p *RuntimeLegacyProviderConfigUsageList)`,
				},
			},
		},
		"Provider config usage list legacy": {
			args: args{
				generate: GenerateProviderConfigUsageListLegacy,
				filename: "zz_generated.legacy_pcu_list.go",
			},
			want: want{
				contains: []string{
					`func (p *RuntimeLegacyProviderConfigUsageList) GetItems() []resource.ProviderConfigUsage {`,
				},
				notContains: []string{
					`func (p *CoreLegacyProviderConfigUsageList)`,
					`func (p *RuntimeModernProviderConfigUsageList)`,
				},
			},
		},
		"Provider config usage list modern core": {
			args: args{
				generate: GenerateProviderConfigUsageListModernCore,
				filename: "zz_generated.modern_core_pcu_list.go",
			},
			want: want{
				contains: []string{
					`func (p *CoreModernProviderConfigUsageList) GetItems() []resource.ProviderConfigUsage {`,
				},
				notContains: []string{
					`func (p *CoreLegacyProviderConfigUsageList)`,
				},
			},
		},
		"Provider config usage list legacy core": {
			args: args{
				generate: GenerateProviderConfigUsageListLegacyCore,
				filename: "zz_generated.legacy_core_pcu_list.go",
			},
			want: want{
				contains: []string{
					`func (p *CoreLegacyProviderConfigUsageList) GetItems() []resource.ProviderConfigUsage {`,
				},
				notContains: []string{
					`func (p *CoreModernProviderConfigUsageList)`,
				},
			},
		},
		"Provider config": {
			args: args{
				generate: GenerateProviderConfig,
				filename: "zz_generated.pc.go",
			},
			want: want{
				contains: []string{
					`func (p *LegacyProviderConfig) SetConditions(c ...xpv1.Condition) {`,
					`func (p *LegacyProviderConfig) GetCondition(ct xpv1.ConditionType) xpv1.Condition {`,
				},
				notContains: []string{
					`func (p *CoreProviderConfig)`,
				},
			},
		},
		"Provider config core": {
			args: args{
				generate: GenerateProviderConfigCore,
				filename: "zz_generated.core_pc.go",
			},
			want: want{
				contains: []string{
					`func (p *CoreProviderConfig) SetUsers(i int64) {`,
					`func (p *CoreProviderConfig) GetUsers() int64 {`,
					`func (p *CoreProviderConfig) SetConditions(c ...xpv2.Condition) {`,
					`func (p *CoreProviderConfig) GetCondition(ct xpv2.ConditionType) xpv2.Condition {`,
				},
				notContains: []string{
					`func (p *CoreProviderConfig) GetCondition(ct xpv1.ConditionType) xpv1.Condition {`,
					`func (p *LegacyProviderConfig)`,
				},
			},
		},
		"References modern": {
			args: args{
				generate: GenerateReferencesModern,
				filename: "zz_generated.modern_refs.go",
			},
			want: want{
				contains: []string{
					`reference.NewAPINamespacedResolver(c, mg)`,
					`var rsp reference.NamespacedResolutionResponse`,
					`func (mg *ModernResource) ResolveReferences(ctx context.Context, c client.Reader) error {`,
				},
			},
		},
		"References legacy": {
			args: args{
				generate: GenerateReferencesLegacy,
				filename: "zz_generated.legacy_refs.go",
			},
			want: want{
				contains: []string{
					`reference.NewAPIResolver(c, mg)`,
					`var rsp reference.ResolutionResponse`,
					`reference.ResolutionRequest{`,
					`func (mg *LegacyResource) ResolveReferences(ctx context.Context, c client.Reader) error {`,
				},
				notContains: []string{
					`reference.NewAPINamespacedResolver(c, mg)`,
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			pkg := loadFixturePackage(t)
			got := generateOutput(t, pkg, tc.args.filename, tc.args.generate)

			for _, want := range tc.want.contains {
				if !strings.Contains(got, want) {
					t.Fatalf("generated output missing snippet %q\n%s", want, got)
				}
			}

			for _, unwanted := range tc.want.notContains {
				if strings.Contains(got, unwanted) {
					t.Fatalf("generated output unexpectedly contained snippet %q\n%s", unwanted, got)
				}
			}
		})
	}
}

func loadFixturePackage(t *testing.T) *packages.Package {
	t.Helper()

	exported := packagestest.Export(t, packagestest.Modules, []packagestest.Module{
		{
			Name: "golang.org/fake",
			Files: map[string]any{
				"v1alpha1/model.go": angryjetFixtureSource,
			},
		},
		{
			Name: "k8s.io/apimachinery",
			Files: map[string]any{
				"pkg/apis/meta/v1/meta.go": metaV1FixtureSource,
			},
		},
		{
			Name: "github.com/crossplane/crossplane",
			Files: map[string]any{
				"apis/v2/core/v2/types.go":          coreV2FixtureSource,
				"apis/v2/core/v2/providerconfig.go": coreV2ProviderConfigFixtureSource,
			},
		},
		{
			Name: "github.com/crossplane/crossplane-runtime/v2",
			Files: map[string]any{
				"apis/common/v1/types.go": runtimeV1FixtureSource,
				"apis/common/v2/types.go": runtimeV2FixtureSource,
			},
		},
	})
	t.Cleanup(exported.Cleanup)

	exported.Config.Mode = LoadMode
	pkgs, err := packages.Load(exported.Config, fmt.Sprintf("file=%s", exported.File("golang.org/fake", "v1alpha1/model.go")))
	if err != nil {
		t.Fatal(err)
	}
	if len(pkgs) != 1 {
		t.Fatalf("packages.Load() returned %d packages, want 1", len(pkgs))
	}
	for _, err := range pkgs[0].Errors {
		t.Fatal(err)
	}

	return pkgs[0]
}

func generateOutput(t *testing.T, pkg *packages.Package, filename string, generate generatorFunc) string {
	t.Helper()

	if err := generate(filename, "", pkg); err != nil {
		t.Fatal(err)
	}

	b, err := os.ReadFile(filepath.Join(filepath.Dir(pkg.GoFiles[0]), filename))
	if err != nil {
		t.Fatal(err)
	}

	return string(b)
}
