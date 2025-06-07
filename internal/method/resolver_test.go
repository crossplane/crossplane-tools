/*
Copyright 2021 The Crossplane Authors.

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

package method

import (
	"fmt"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/packages/packagestest"

	"github.com/crossplane/crossplane-tools/internal/comments"
	xptypes "github.com/crossplane/crossplane-tools/internal/types"
)

const (
	source = `
package v1alpha1

type ModelParameters struct {
	// +crossplane:generate:reference:type=Apigatewayv2Api
	APIID string

	// +crossplane:generate:reference:type=SecurityGroup
	SecurityGroupID *string

	// +crossplane:generate:reference:type=github.com/crossplane/provider-aws/apis/identity/v1beta1.IAM
	// +crossplane:generate:reference:extractor=github.com/crossplane/provider-aws/apis/identity/v1beta1.IAMRoleARN()
	IAMRoleARN *string

	// +crossplane:generate:reference:type=github.com/crossplane/provider-aws/apis/identity/v1beta1.IAM
	// +crossplane:generate:reference:extractor=github.com/crossplane/provider-aws/apis/identity/v1beta1.IAMRoleARN("a.b.c")
	NestedTargetWithPath *string

	// +crossplane:generate:reference:type=github.com/crossplane/provider-aws/apis/identity/v1beta1.IAM
	// +crossplane:generate:reference:extractor=IAMRoleARN("a.b.c")
	NestedTargetNoPath *string

	// +crossplane:generate:reference:type=github.com/crossplane/provider-aws/apis/identity/v1beta1.IAM
	// +crossplane:generate:reference:extractor=IAMRoleARN()
	NoArgNoPath *string

	Network *NetworkSpec

	OtherSetting []OtherSpec

	// +crossplane:generate:reference:type=Subnet
	// +crossplane:generate:reference:refFieldName=SubnetIDRefs
	// +crossplane:generate:reference:selectorFieldName=SubnetIDSelector
	SubnetIDs []string

	// +crossplane:generate:reference:type=RouteTable
	RouteTableIDs []*string

	UnrelatedField string

	// +crossplane:generate:reference:type=golang.org/fake/v1alpha1.Configuration
	// +crossplane:generate:reference:extractor=golang.org/fake/v1alpha1.Configuration()
	CustomConfiguration *Configuration

	// +crossplane:generate:reference:type=github.com/crossplane/provider-aws/apis/identity/v1beta1.IAM
	// +crossplane:generate:reference:extractor=Count()
	Count *float64
}

type NetworkSpec struct {
	// +crossplane:generate:reference:type=github.com/crossplane/provider-aws/apis/ec2/v1beta1.VPC
	VPCID string
}

type OtherSpec struct {
	// +crossplane:generate:reference:type=Cluster
	OtherID string
}

type ModelSpec struct {
	ForProvider       ModelParameters
}

type ModelObservation struct {}

type ModelStatus struct {
	AtProvider          ModelObservation
}

type Model struct {
	Spec              ModelSpec
	Status            ModelStatus
}
`
	generated = `package v1alpha1

import (
	"context"
	client "example.org/client"
	reference "example.org/reference"
	v1beta11 "github.com/crossplane/provider-aws/apis/ec2/v1beta1"
	v1beta1 "github.com/crossplane/provider-aws/apis/identity/v1beta1"
	errors "github.com/pkg/errors"
)

// ResolveReferences of this Model.
func (mg *Model) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, mg)

	var rsp reference.ResolutionResponse
	var mrsp reference.MultiResolutionResponse
	var err error

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: mg.Spec.ForProvider.APIID,
		Extract:      reference.ExternalName(),
		Namespace:    mg.GetNamespace(),
		Reference:    mg.Spec.ForProvider.APIIDRef,
		Selector:     mg.Spec.ForProvider.APIIDSelector,
		To: reference.To{
			List:    &Apigatewayv2ApiList{},
			Managed: &Apigatewayv2Api{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.APIID")
	}
	mg.Spec.ForProvider.APIID = rsp.ResolvedValue
	mg.Spec.ForProvider.APIIDRef = rsp.ResolvedReference

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.SecurityGroupID),
		Extract:      reference.ExternalName(),
		Namespace:    mg.GetNamespace(),
		Reference:    mg.Spec.ForProvider.SecurityGroupIDRef,
		Selector:     mg.Spec.ForProvider.SecurityGroupIDSelector,
		To: reference.To{
			List:    &SecurityGroupList{},
			Managed: &SecurityGroup{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.SecurityGroupID")
	}
	mg.Spec.ForProvider.SecurityGroupID = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.SecurityGroupIDRef = rsp.ResolvedReference

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.IAMRoleARN),
		Extract:      v1beta1.IAMRoleARN(),
		Namespace:    mg.GetNamespace(),
		Reference:    mg.Spec.ForProvider.IAMRoleARNRef,
		Selector:     mg.Spec.ForProvider.IAMRoleARNSelector,
		To: reference.To{
			List:    &v1beta1.IAMList{},
			Managed: &v1beta1.IAM{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.IAMRoleARN")
	}
	mg.Spec.ForProvider.IAMRoleARN = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.IAMRoleARNRef = rsp.ResolvedReference

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.NestedTargetWithPath),
		Extract:      v1beta1.IAMRoleARN("a.b.c"),
		Namespace:    mg.GetNamespace(),
		Reference:    mg.Spec.ForProvider.NestedTargetWithPathRef,
		Selector:     mg.Spec.ForProvider.NestedTargetWithPathSelector,
		To: reference.To{
			List:    &v1beta1.IAMList{},
			Managed: &v1beta1.IAM{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.NestedTargetWithPath")
	}
	mg.Spec.ForProvider.NestedTargetWithPath = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.NestedTargetWithPathRef = rsp.ResolvedReference

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.NestedTargetNoPath),
		Extract:      IAMRoleARN("a.b.c"),
		Namespace:    mg.GetNamespace(),
		Reference:    mg.Spec.ForProvider.NestedTargetNoPathRef,
		Selector:     mg.Spec.ForProvider.NestedTargetNoPathSelector,
		To: reference.To{
			List:    &v1beta1.IAMList{},
			Managed: &v1beta1.IAM{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.NestedTargetNoPath")
	}
	mg.Spec.ForProvider.NestedTargetNoPath = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.NestedTargetNoPathRef = rsp.ResolvedReference

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.NoArgNoPath),
		Extract:      IAMRoleARN(),
		Namespace:    mg.GetNamespace(),
		Reference:    mg.Spec.ForProvider.NoArgNoPathRef,
		Selector:     mg.Spec.ForProvider.NoArgNoPathSelector,
		To: reference.To{
			List:    &v1beta1.IAMList{},
			Managed: &v1beta1.IAM{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.NoArgNoPath")
	}
	mg.Spec.ForProvider.NoArgNoPath = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.NoArgNoPathRef = rsp.ResolvedReference

	if mg.Spec.ForProvider.Network != nil {
		rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
			CurrentValue: mg.Spec.ForProvider.Network.VPCID,
			Extract:      reference.ExternalName(),
			Namespace:    mg.GetNamespace(),
			Reference:    mg.Spec.ForProvider.Network.VPCIDRef,
			Selector:     mg.Spec.ForProvider.Network.VPCIDSelector,
			To: reference.To{
				List:    &v1beta11.VPCList{},
				Managed: &v1beta11.VPC{},
			},
		})
		if err != nil {
			return errors.Wrap(err, "mg.Spec.ForProvider.Network.VPCID")
		}
		mg.Spec.ForProvider.Network.VPCID = rsp.ResolvedValue
		mg.Spec.ForProvider.Network.VPCIDRef = rsp.ResolvedReference

	}
	for i3 := 0; i3 < len(mg.Spec.ForProvider.OtherSetting); i3++ {
		rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
			CurrentValue: mg.Spec.ForProvider.OtherSetting[i3].OtherID,
			Extract:      reference.ExternalName(),
			Namespace:    mg.GetNamespace(),
			Reference:    mg.Spec.ForProvider.OtherSetting[i3].OtherIDRef,
			Selector:     mg.Spec.ForProvider.OtherSetting[i3].OtherIDSelector,
			To: reference.To{
				List:    &ClusterList{},
				Managed: &Cluster{},
			},
		})
		if err != nil {
			return errors.Wrap(err, "mg.Spec.ForProvider.OtherSetting[i3].OtherID")
		}
		mg.Spec.ForProvider.OtherSetting[i3].OtherID = rsp.ResolvedValue
		mg.Spec.ForProvider.OtherSetting[i3].OtherIDRef = rsp.ResolvedReference

	}
	mrsp, err = r.ResolveMultiple(ctx, reference.MultiResolutionRequest{
		CurrentValues: mg.Spec.ForProvider.SubnetIDs,
		Extract:       reference.ExternalName(),
		Namespace:     mg.GetNamespace(),
		References:    mg.Spec.ForProvider.SubnetIDRefs,
		Selector:      mg.Spec.ForProvider.SubnetIDSelector,
		To: reference.To{
			List:    &SubnetList{},
			Managed: &Subnet{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.SubnetIDs")
	}
	mg.Spec.ForProvider.SubnetIDs = mrsp.ResolvedValues
	mg.Spec.ForProvider.SubnetIDRefs = mrsp.ResolvedReferences

	mrsp, err = r.ResolveMultiple(ctx, reference.MultiResolutionRequest{
		CurrentValues: reference.FromPtrValues(mg.Spec.ForProvider.RouteTableIDs),
		Extract:       reference.ExternalName(),
		Namespace:     mg.GetNamespace(),
		References:    mg.Spec.ForProvider.RouteTableIDsRefs,
		Selector:      mg.Spec.ForProvider.RouteTableIDsSelector,
		To: reference.To{
			List:    &RouteTableList{},
			Managed: &RouteTable{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.RouteTableIDs")
	}
	mg.Spec.ForProvider.RouteTableIDs = reference.ToPtrValues(mrsp.ResolvedValues)
	mg.Spec.ForProvider.RouteTableIDsRefs = mrsp.ResolvedReferences

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.CustomConfiguration),
		Extract:      Configuration(),
		Namespace:    mg.GetNamespace(),
		Reference:    mg.Spec.ForProvider.CustomConfigurationRef,
		Selector:     mg.Spec.ForProvider.CustomConfigurationSelector,
		To: reference.To{
			List:    &ConfigurationList{},
			Managed: &Configuration{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.CustomConfiguration")
	}
	mg.Spec.ForProvider.CustomConfiguration = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CustomConfigurationRef = rsp.ResolvedReference

	rsp, err = r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromFloatPtrValue(mg.Spec.ForProvider.Count),
		Extract:      Count(),
		Namespace:    mg.GetNamespace(),
		Reference:    mg.Spec.ForProvider.CountRef,
		Selector:     mg.Spec.ForProvider.CountSelector,
		To: reference.To{
			List:    &v1beta1.IAMList{},
			Managed: &v1beta1.IAM{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "mg.Spec.ForProvider.Count")
	}
	mg.Spec.ForProvider.Count = reference.ToFloatPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.CountRef = rsp.ResolvedReference

	return nil
}
`
)

func TestNewResolveReferences(t *testing.T) {
	exported := packagestest.Export(t, packagestest.Modules, []packagestest.Module{{
		Name: "golang.org/fake",
		Files: map[string]any{
			"v1alpha1/model.go": source,
		},
	}})
	defer exported.Cleanup()
	exported.Config.Mode = packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedDeps | packages.NeedTypes | packages.NeedSyntax
	pkgs, err := packages.Load(exported.Config, fmt.Sprintf("file=%s", exported.File("golang.org/fake", "v1alpha1/model.go")))
	if err != nil {
		t.Error(err)
	}
	f := jen.NewFilePath("golang.org/fake/v1alpha1")
	NewResolveReferences(xptypes.NewTraverser(comments.In(pkgs[0])), "mg", "example.org/client", "example.org/reference")(f, pkgs[0].Types.Scope().Lookup("Model"))
	if diff := cmp.Diff(generated, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewResolveReferences(): -want, +got\n%s", diff)
	}
}
