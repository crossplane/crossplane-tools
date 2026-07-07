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

type ModernProviderConfigUsage struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	xpv2.TypedProviderConfigUsage
}

type ModernProviderConfigUsageList struct {
	metav1.TypeMeta
	Items []ModernProviderConfigUsage
}

type ModernProviderCredentials struct {
	Source xpv2.CredentialsSource
	xpv2.CommonCredentialSelectors
}

type ModernProviderConfigSpec struct {
	xprv1.ProviderConfigSpec
	Credentials ModernProviderCredentials
}

// A ModernProviderConfigStatus reflects the observed state of a ProviderConfig.
type ModernProviderConfigStatus struct {
	xpv2.ProviderConfigStatus
}

type ModernProviderConfig struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   ModernProviderConfigSpec
	Status ModernProviderConfigStatus
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

type ResourceStatus struct{}

type Condition struct{}

type ConditionType string

type User struct{}

type ProviderConfigSpec struct{}

type ProviderConfigStatus struct{}

func (p *ProviderConfigStatus) SetConditions(c ...Condition) {}

func (p *ProviderConfigStatus) GetCondition(ct ConditionType) Condition { return Condition{} }

func (p *ProviderConfigStatus) SetUsers(u ...User) {}

func (p *ProviderConfigStatus) GetUsers() []User { return nil }
`

	runtimeV2FixtureSource = `package v2

type ManagedResourceSpec struct{}
`

	coreV2ProviderConfigFixtureSource = `package v2

type ProviderConfigStatus struct{}

type ProviderConfigReference struct{}

type TypedReference struct{}

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
		"Managed namespaced": {
			args: args{
				generate: GenerateManagedNamespaced,
				filename: "zz_generated.namespaced.go",
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
		"Managed cluster": {
			args: args{
				generate: GenerateManagedCluster,
				filename: "zz_generated.cluster.go",
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
		"Managed list namespaced": {
			args: args{
				generate: GenerateManagedNamespacedList,
				filename: "zz_generated.namespaced_list.go",
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
		"Managed list cluster": {
			args: args{
				generate: GenerateManagedClusterList,
				filename: "zz_generated.cluster_list.go",
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
		"References namespaced": {
			args: args{
				generate: GenerateReferencesNamespaced,
				filename: "zz_generated.namespaced_refs.go",
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
		"References cluster": {
			args: args{
				generate: GenerateReferencesCluster,
				filename: "zz_generated.cluster_refs.go",
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
		"Provider config usage modern": {
			args: args{
				generate: GenerateProviderConfigUsageModern,
				filename: "zz_generated.modern_pcu.go",
			},
			want: want{
				contains: []string{
					`func (p *ModernProviderConfigUsage) SetProviderConfigReference(r xpv2.ProviderConfigReference) {`,
					`func (p *ModernProviderConfigUsage) GetProviderConfigReference() xpv2.ProviderConfigReference {`,
					`func (p *ModernProviderConfigUsage) SetResourceReference(r xpv2.TypedReference) {`,
					`func (p *ModernProviderConfigUsage) GetResourceReference() xpv2.TypedReference {`,
				},
				notContains: []string{
					`func (p *ModernProviderConfigUsage) SetProviderConfigReference(r xpv1.ProviderConfigReference) {`,
					`func (p *ModernProviderConfigUsage) GetResourceReference() xpv1.TypedReference {`,
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
					`func (p *ModernProviderConfigUsageList) GetItems() []resource.ProviderConfigUsage {`,
				},
			},
		},
		"Provider config legacy": {
			args: args{
				generate: GenerateProviderConfigLegacy,
				filename: "zz_generated.legacy_pc.go",
			},
			want: want{
				contains: []string{
					`func (p *LegacyProviderConfig) SetConditions(c ...xpv1.Condition) {`,
					`func (p *LegacyProviderConfig) GetCondition(ct xpv1.ConditionType) xpv1.Condition {`,
				},
			},
		},
		"Provider config modern status and credential migration": {
			args: args{
				generate: GenerateProviderConfigModern,
				filename: "zz_generated.modern_pc.go",
			},
			want: want{
				contains: []string{
					`func (p *ModernProviderConfig) SetUsers(i int64) {`,
					`func (p *ModernProviderConfig) GetUsers() int64 {`,
					`func (p *ModernProviderConfig) SetConditions(c ...xpv2.Condition) {`,
					`func (p *ModernProviderConfig) GetCondition(ct xpv2.ConditionType) xpv2.Condition {`,
				},
				notContains: []string{
					`func (p *ModernProviderConfig) GetCondition(ct xpv1.ConditionType) xpv1.Condition {`,
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
