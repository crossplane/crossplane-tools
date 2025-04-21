/*
Copyright 2025 The Crossplane Authors.

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

package convert

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestToAndFromPtr(t *testing.T) {
	cases := map[string]struct {
		want string
	}{
		"Zero":    {want: ""},
		"NonZero": {want: "pointy"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := FromPtrValue(ToPtrValue(tc.want))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("FromPtrValue(ToPtrValue(%s): -want, +got: %s", tc.want, diff)
			}
		})
	}
}

func TestToAndFromFloatPtr(t *testing.T) {
	cases := map[string]struct {
		want string
	}{
		"Zero":    {want: ""},
		"NonZero": {want: "1123581321"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := FromFloatPtrValue(ToFloatPtrValue(tc.want))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("FromPtrValue(ToPtrValue(%s): -want, +got: %s", tc.want, diff)
			}
		})
	}
}

func TestToAndFromPtrValues(t *testing.T) {
	cases := map[string]struct {
		want []string
	}{
		"Nil":      {want: []string{}},
		"Zero":     {want: []string{""}},
		"NonZero":  {want: []string{"pointy"}},
		"Multiple": {want: []string{"pointy", "pointers"}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := FromPtrValues(ToPtrValues(tc.want))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("FromPtrValues(ToPtrValues(%s): -want, +got: %s", tc.want, diff)
			}
		})
	}
}

func TestToAndFromFloatPtrValues(t *testing.T) {
	cases := map[string]struct {
		want []string
	}{
		"Nil":      {want: []string{}},
		"Zero":     {want: []string{""}},
		"NonZero":  {want: []string{"1123581321"}},
		"Multiple": {want: []string{"1123581321", "1234567890"}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := FromFloatPtrValues(ToFloatPtrValues(tc.want))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("FromPtrValues(ToPtrValues(%s): -want, +got: %s", tc.want, diff)
			}
		})
	}
}

func TestToAndFromIntPtrValues(t *testing.T) {
	cases := map[string]struct {
		want []string
	}{
		"Nil":      {want: []string{}},
		"Zero":     {want: []string{""}},
		"NonZero":  {want: []string{"1123581321"}},
		"Multiple": {want: []string{"1123581321", "1234567890"}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := FromIntPtrValues(ToIntPtrValues(tc.want))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("FromIntPtrValues(ToIntPtrValues(%s): -want, +got: %s", tc.want, diff)
			}
		})
	}
}
