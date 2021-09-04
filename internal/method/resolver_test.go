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
)

func TestNewResolveReferences(t *testing.T) {
	want := `package pkg

import (
	"context"
	client "example.org/client"
	reference "example.org/reference"
	errors "github.com/pkg/errors"
)

// ResolveReferences of this Type.
func (t *Type) ResolveReferences(ctx context.Context, c client.Reader) {
	r := reference.NewAPIResolver(c, t)

	resp, err := r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: t.Spec.ForProvider.ApiId,
		Extract:      reference.ExternalName(),
		Reference:    t.Spec.ForProvider.ApiIdRef,
		Selector:     t.Spec.ForProvider.ApiIdSelector,
		To: reference.To{
			List:    &Apigatewayv2ApiList{},
			Managed: &Apigatewayv2Api{},
		},
	})
	if err != nil {
		return errors.Wrap(err, "Spec.ForProvider.ApiId")
	}

	t.Spec.ForProvider.ApiId = resp.ResolvedValue
	t.Spec.ForProvider.ApiIdRef = resp.ResolvedReference

	return nil
}
`
	f := jen.NewFile("pkg")
	NewResolveReferences("t", "example.org/client", "example.org/reference")(f, MockObject{Named: "Type"})
	if diff := cmp.Diff(want, fmt.Sprintf("%#v", f)); diff != "" {
		t.Errorf("NewResolveReferences(): -want, +got\n%s", diff)
	}
}
